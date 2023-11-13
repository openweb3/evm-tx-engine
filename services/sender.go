package services

import (
	"github.com/openweb3/evm-tx-engine/models"
	"github.com/openweb3/evm-tx-engine/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Fetch a batch of transaction from signed queue, and send to chain
// if the transaction is not a cancel transaction, make sure the related task is not in cancelling status

// temp impl
// don't actually send but mark as pending
func StartSenderRound(db *gorm.DB) {
	var txs []models.ChainTransaction

	err := db.Model(&models.ChainTransaction{}).Joins("Field").Where("tx_status = ? AND raw IS not NULL", utils.TxInternalSigned).Find(&txs).Error

	if err != nil {
		logrus.WithError(err).Error("Failed to get signed transactions")
		return
	}

	if len(txs) == 0 {
		return
	}

	for i := range txs {
		txs[i].TxStatus = utils.TxPoolPending
	}
	err = db.Save(&txs).Error
	if err != nil {
		logrus.WithError(err).Error("Failed to update signed transactions")
		return
	}
	logrus.WithField("service", "sender").Infof("batch sent %d transaction(s)", len(txs))

}

// temp impl
// don't actually send but mark as pending
func StartSenderWorkerRound(ctx *QueueContext, maxBatchSize int) error {
	var txs = ctx.SendingQueue.MustDequeBatch(maxBatchSize)
	if len(txs) == 0 {
		return nil
	}
	backupTxs := backupChainTransactions(txs)

	err := func() error {
		// TODO: send transactions to node here
		for _, tx := range txs {
			tx.TxStatus = utils.TxPoolPending
		}
		return nil
	}()

	// processes sending error
	if err != nil {
		txs = backupTxs
		logrus.WithError(err).Error("Failed to send transactions")
		for _, tx := range txs {
			ctx.ErrQueue.MustEnqueWithLog(*tx, "Sender", "error sending transaction")
		}
		return err
	}

	err = models.SaveWithRetry(ctx.Db, &txs)
	if err != nil {
		logrus.WithError(err).Error("failed to save transactions")
		return err
	}

	for _, tx := range txs {
		ctx.PoolOrChainQueue.MustEnqueWithLog(*tx, "Sender", "transaction sent")
	}
	logrus.WithField("service", "sender").Infof("batch sent %d transaction(s)", len(txs))
	return nil
}
