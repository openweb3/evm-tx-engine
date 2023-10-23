package services

import (
	"github.com/openweb3/evm-tx-engine/models"
	"github.com/openweb3/evm-tx-engine/utils"
	"gorm.io/gorm"
)

// Scan GasEnoughQueue, and modify the gas field, then change the transaction to constructed (all fields are supposed to complete)

func StartPriceManageRound(db *gorm.DB) {
	var chainTransactions []models.ChainTransaction
	db.Preload("Fields").Where("tx_status = ? AND nonce IS not NULL", utils.TxInternalGasEnoughQueue).Find(&chainTransactions)
	// ignore the case that gas price was already filled in the simple implememtation
	for _, chainTransaction := range chainTransactions {
		chainTransaction.Field.GasPrice = 5
		chainTransaction.TxStatus = utils.TxInternalConstructed
	}
	db.Save(&chainTransactions)
}
