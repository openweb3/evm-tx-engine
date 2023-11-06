package services

import (
	"github.com/openweb3/evm-tx-engine/models"
	"github.com/openweb3/evm-tx-engine/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// estimates gas consumption
func StartExecutionSimulationRound(db *gorm.DB) {
	var txs []models.ChainTransaction

	db.Model(&models.ChainTransaction{}).Preload("Field").Where("tx_status = ?", utils.TxInternalTargetQueue).Find(&txs)

	if len(txs) == 0 {
		return
	}

	for _, tx := range txs {
		if tx.Field.GasLimit != 0 {
			continue
		}
		tx.Field.GasLimit = 10000000
		err := db.Save(&tx.Field).Error
		if err != nil {
			logrus.WithField("service", "execution simulator").WithError(err).Error("failed to save transaction gas limit field")
			continue
		}
		logrus.WithField("service", "execution simulator").WithField("txId", tx.ID).Info("gas limit set for transaction")
	}
}

func StartExecutionSimulationWorkerRound(ctx *QueueContext, maxBatchSize int) error {
	txs := ctx.TargetQueue.MustDequeBatch(maxBatchSize)

	for _, tx := range *txs {
		if tx.Field.GasLimit != 0 {
			continue
		}
		tx.Field.GasLimit = 10000000
		ctx.BalanceCheckingQueue.MustEnqueWithLog(tx, "ExecutionSimulator", "gas limit set")
	}
	return nil
}
