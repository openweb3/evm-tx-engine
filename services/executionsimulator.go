package services

import (
	"github.com/openweb3/evm-tx-engine/models"
	"github.com/openweb3/evm-tx-engine/utils"
	"gorm.io/gorm"
)

// estimates gas consumption

func StartExecutionSimulationRound(db *gorm.DB) {
	var txs []models.ChainTransaction

	db.Model(&models.ChainTransaction{}).Preload("Field").Where("tx_status = ?", utils.TxInternalTargetQueue).Find(&txs)

	for _, tx := range txs {
		tx.Field.GasLimit = 10000000
	}
	db.Save(&txs)
}
