package services

import (
	"github.com/openweb3/evm-tx-engine/accountadapter"
	"github.com/openweb3/evm-tx-engine/models"
	"github.com/openweb3/evm-tx-engine/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// pop a batch of transactions from constructed, and send them to the signing service. Modify the `Raw` field and modify the status field
func StartSigningRound(db *gorm.DB) {
	var txs []models.ChainTransaction
	db.Where("tx_status = ? AND raw IS NULL", utils.TxInternalConstructed).Find(&txs)
	// will modify txs
	if len(txs) == 0 {
		return
	}
	err := batchSign(db, &txs)
	if err != nil {
		// do something
		logrus.WithField("service", "signing").WithError(err).Error("error signing transaction")
		return
	}
	for i := range txs {
		txs[i].TxStatus = utils.TxInternalSigned
	}
	err = db.Save(&txs).Error
	if err != nil {
		logrus.WithField("service", "signing").WithError(err).Error("Error saving transaction")
		return
	}
	logrus.WithField("service", "signing").Infof("batch signed %d transaction(s)", len(txs))

}

func batchSign(db *gorm.DB, txs *[]models.ChainTransaction) error {
	for i := range *txs {
		txReadyToSign, err := (*txs)[i].PrepareTransactionToSign()
		if err != nil {
			return err
		}
		fromAccount, err := (*txs)[i].GetTransactionFrom(db)
		if err != nil {
			return err
		}
		signer := accountadapter.Signer

		signed, err := signer.SignTransaction(fromAccount.Address, txReadyToSign)
		if err != nil {
			return err
		}
		(*txs)[i].Raw = &signed
	}
	return nil
}

// TODO: batch save
func StartSigningWorkerRound(ctx *QueueContext, maxBatchSize int) error {
	txs := ctx.SigningQueue.MustDequeBatch(maxBatchSize)
	if len(*txs) == 0 {
		return nil
	}

	backupTxs := *txs

	err := func() error {
		err := batchSign(ctx.Db, txs)
		if err != nil {
			// do something
			logrus.WithField("service", "signing").WithError(err).Error("error signing transactions")
			return err
		}
		for i := range *txs {
			(*txs)[i].TxStatus = utils.TxInternalSigned
		}
		err = ctx.Db.Save(&txs).Error
		return err
	}()

	if err != nil {
		*txs = backupTxs
		logrus.WithField("service", "signing").WithError(err).Error("Error saving transaction")
		for _, tx := range *txs {
			ctx.ErrQueue.MustEnqueWithLog(tx, "SigningService", "error saving transaction")
		}
		return err
	}
	for _, tx := range *txs {
		ctx.SendingQueue.MustEnqueWithLog(tx, "SigningService", "transaction signed")
	}
	return nil
}
