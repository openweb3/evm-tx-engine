package services

import (
	"errors"

	"github.com/openweb3/evm-tx-engine/models"
	"github.com/openweb3/evm-tx-engine/types/code"
	"gorm.io/gorm"
)

func StartNonceManageWorkerRound(ctx *QueueContext, maxBatchSize int) error {
	var txs = ctx.NonceManagingQueue.MustDequeBatch(maxBatchSize)
	if len(txs) == 0 {
		return nil
	}
	for _, tx := range txs {
		err := SetTransactionNonce(ctx.Db, tx)
		if err != nil {
			ctx.ErrQueue.MustEnqueWithLog(*tx, "NonceManager", "set nonce error")
		}
		ctx.SigningQueue.MustEnqueWithLog(*tx, "NonceManager", "tx nonce allocated")
	}
	return nil
}

// NOTE: never directly change the value pointer points to but provide a new pointer!
// wrong: *tx.a = newVal
// right: b := newVal; tx.a = b
// SetTransactionNonce is a function that sets the nonce of a transaction
// It takes a pointer to a gorm.DB and a pointer to a ChainTransaction as parameters
// It returns an error if there is one
func SetTransactionNonce(db *gorm.DB, tx *models.ChainTransaction) error {
	// start a transaction
	dbTransaction := db.Begin()
	// back up the transaction
	backupTx := *tx
	// try to set the nonce
	err := func() error {
		// should be nil
		if tx.Field.Nonce != nil {
			return errors.New("nonce already allocated")
		}
		// get the account from the transaction
		fromAccount, err := tx.GetTransactionFrom(dbTransaction)
		if err != nil {
			return err
		}
		// set the new nonce
		newNonce := fromAccount.InternalNonce
		tx.Field.Nonce = &newNonce
		tx.TxStatus = code.TxInternalSigning
		// update the fromAccount internal nonce
		// TODO: should lock fromAccount if multiple workers are working
		fromAccount.InternalNonce = fromAccount.InternalNonce + 1

		err = models.SaveWithRetry(dbTransaction, &fromAccount)
		if err != nil {
			return err
		}

		err = models.SaveWithRetry(dbTransaction, &tx)
		if err != nil {
			return err
		}
		return nil
	}()
	// if there is an error, rollback the transaction and restore the tx
	if err != nil {
		*tx = backupTx // restore tx
		dbTransaction.Rollback()
		return err
	}
	// commit the transaction
	dbTransaction.Commit()
	return nil
}
