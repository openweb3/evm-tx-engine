package services

import (
	"errors"

	"github.com/openweb3/evm-tx-engine/models"
	"github.com/openweb3/evm-tx-engine/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Scan GasEnoughQueue, and modify the nonce field
// no lock in the first impl
// TODO: lock for account internal nonce
// As nonce is quite complicated, don't use batch operation here
func StartNonceManageRound(db *gorm.DB) {
	var txs []models.ChainTransaction
	db.Joins("Field", db.Where("nonce IS NULL")).Find(&txs, "tx_status = ?", utils.TxInternalGasEnoughQueue)
	for _, tx := range txs {
		dbTransaction := db.Begin()
		fromAccount, err := tx.GetTransactionFrom(dbTransaction)
		if err != nil {
			dbTransaction.Rollback()
			continue
		}
		// should be nil because we filter it
		if tx.Field.Nonce != nil {
			dbTransaction.Rollback()
			continue
		}
		newNonce := fromAccount.InternalNonce
		tx.Field.Nonce = &newNonce
		err = dbTransaction.Save(&tx.Field).Error
		if err != nil {
			dbTransaction.Rollback()
			continue
		}
		fromAccount.InternalNonce = fromAccount.InternalNonce + 1
		err = dbTransaction.Save(&tx).Error
		if err != nil {
			dbTransaction.Rollback()
			continue
		}
		dbTransaction.Commit()
		logrus.WithField("service", "nonce manage").WithField("txId", tx.ID).Info("nonce set for transaction")
	}
}

func StartNonceManageWorkderRound(ctx *QueueContext, maxBatchSize int) error {
	var txs = ctx.NonceManagingQueue.MustDequeBatch(maxBatchSize)
	if len(*txs) == 0 {
		return nil
	}
	for _, tx := range *txs {
		err := SetTransactionNonce(ctx.Db, &tx)
		if err != nil {
			ctx.ErrQueue.MustEnqueWithLog(tx, "NonceManager", "set nonce error")
		}
		ctx.SigningQueue.MustEnqueWithLog(tx, "NonceManager", "tx nonce allocated")
	}
	return nil
}

// NOTE: never directly change the value pointer points to but provide a new pointer!
// wrong: *tx.a = newVal
// right: b := newVal; tx.a = b
func SetTransactionNonce(db *gorm.DB, tx *models.ChainTransaction) error {
	dbTransaction := db.Begin()
	backupTx := *tx
	err := func() error {
		// should be nil
		if tx.Field.Nonce != nil {
			return errors.New("nonce already allocated")
		}
		fromAccount, err := tx.GetTransactionFrom(dbTransaction)
		if err != nil {
			return err
		}
		newNonce := fromAccount.InternalNonce
		tx.Field.Nonce = &newNonce
		// save field
		err = dbTransaction.Save(&tx.Field).Error
		if err != nil {
			return err
		}
		// update fromAccount internal nonce
		// TODO: should lock fromAccount if multiple workers are working
		fromAccount.InternalNonce = fromAccount.InternalNonce + 1
		err = db.Save(&fromAccount).Error
		if err != nil {
			return err
		}
		// change txStatus
		tx.TxStatus = utils.TxInternalConstructed
		err = dbTransaction.Save(&tx).Error
		if err != nil {
			return err
		}
		return nil
	}()
	if err != nil {
		*tx = backupTx // restore tx
		dbTransaction.Rollback()
		return err
	}
	dbTransaction.Commit()
	return nil
}
