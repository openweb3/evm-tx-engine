package services

import (
	"github.com/openweb3/evm-tx-engine/models"
	"github.com/openweb3/evm-tx-engine/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Scan GasEnoughQueue, and modify the nonce field
// no lock in the first impl
// TODO: lock for account internal nonce
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
