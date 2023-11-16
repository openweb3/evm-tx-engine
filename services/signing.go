package services

import (
	"github.com/openweb3/evm-tx-engine/accountadapter"
	"github.com/openweb3/evm-tx-engine/models"
	"github.com/openweb3/evm-tx-engine/types/code"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func batchSign(db *gorm.DB, txs []*models.ChainTransaction) error {
	for _, tx := range txs {
		fromAccount, err := tx.GetTransactionFrom(db)
		if err != nil {
			return err
		}
		signer := accountadapter.Signer

		signed, err := signer.SignTransaction(fromAccount.Address, fromAccount.Chain.ChainId, tx.Field)
		if err != nil {
			return err
		}
		tx.Raw = &signed
	}
	return nil
}

// TODO: batch save
func StartSigningWorkerRound(ctx *QueueContext, maxBatchSize int) error {
	txs := ctx.SigningQueue.MustDequeBatch(maxBatchSize)
	if len(txs) == 0 {
		return nil
	}

	backupTxs := backupChainTransactions(txs)

	// sign transaction operations
	err := func() error {
		err := batchSign(ctx.Db, txs)
		if err != nil {
			return err
		}
		for _, tx := range txs {
			tx.TxStatus = code.TxInternalSending
		}
		return nil
	}()

	// recover signing errors
	if err != nil {
		txs = backupTxs
		logrus.WithField("service", "signing").WithError(err).Error("Error signing transaction")
		for _, tx := range txs {
			ctx.ErrQueue.MustEnqueWithLog(*tx, "SigningService", "error saving transaction")
		}
		return err
	}

	// logs the errors
	err = models.SaveWithRetry(ctx.Db, &txs)
	if err != nil {
		logrus.WithError(err).Error("error saving transactions to db: &+v", txs)
		return err
	}

	for _, tx := range txs {
		ctx.SendingQueue.MustEnqueWithLog(*tx, "SigningService", "transaction signed")
	}
	return nil
}
