package services

import (
	"github.com/openweb3/evm-tx-engine/models"
	"github.com/openweb3/evm-tx-engine/utils"
	"gorm.io/gorm"
)

// Fetch a batch of transaction from signed queue, and send to chain
// if the transaction is not a cancel transaction, make sure the related task is not in cancelling status

// temp impl
// don't actually send but mark as pending
func StartSenderRound(db *gorm.DB) {
	var txs []models.ChainTransaction

	db.Model(&models.ChainTransaction{}).Preload("Field").Where("tx_status = ?", utils.TxInternalSigned).Find(&txs)

	for _, tx := range txs {
		tx.TxStatus = utils.TxPoolPending
	}
	db.Save(&txs)
}
