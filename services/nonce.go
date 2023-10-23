package services

import (
	"github.com/openweb3/evm-tx-engine/models"
	"github.com/openweb3/evm-tx-engine/utils"
	"gorm.io/gorm"
)

// Scan GasEnoughQueue, and modify the nonce field
// no lock in the first impl
func StartNonceManageRound(db *gorm.DB) {
	var txs []models.ChainTransaction
	db.Joins("Fields", db.Where("nonce IS not NULL")).Preload("task.from").Find(&txs, "tx_status = ?", utils.TxInternalGasEnoughQueue)
	for _, tx := range txs {
		dbTransaction := db.Begin()
		err := db.Model(&tx).Update("nonce", tx.Task.From.InternalNonce).Error
		if err != nil {
			dbTransaction.Rollback()
			continue
		}
		tx.Task.From.InternalNonce += 1
		db.Save(tx)
		dbTransaction.Commit()
	}
}
