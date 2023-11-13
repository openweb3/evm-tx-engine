package services

import (
	"errors"

	"github.com/openweb3/evm-tx-engine/models"
	"github.com/openweb3/evm-tx-engine/types"
	"github.com/openweb3/evm-tx-engine/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// leave as empty

// poll eth latest/safe/finalized blocknumber (TODO: for each chain)
// Modify LatestBlock, SafeBlock, FinalizedBlock
func StartTaggedBlockNumberUpdateRound(db *gorm.DB) {
	var chains []models.Chain
	db.Find(&chains)
	for _, chain := range chains {
		taggedBlockNumbers, err := utils.W3.GetTaggedBlockNumbers(chain.Name)
		if err != nil {
			logrus.WithError(err).Error("Failed to get tagged block numbers")
			continue
		}
		chain.LatestBlock = taggedBlockNumbers.LatestBlock
		chain.SafeBlock = taggedBlockNumbers.SafeBlock
		chain.FinalizedBlock = taggedBlockNumbers.FinalizedBlock
		err = db.Save(&chain).Error
		if err != nil {
			logrus.WithError(err).Error("Failed to update tagged block numbers")
			continue
		}
		// logrus.WithField("service", "ChainUpdate").Info("Updated tagged block numbers")
	}
	// TODO: panic or something
}

// Update TransactionBlockNumber
// Update Transaction Status
// Update UpdatedAt if the transaction is of insufficient balance (including txstatusinternal transaction) (?)
// the IsStable will also be updated. But note this field might be not the right value now because certain related value is not updated.
func StartTransactionChainStatusUpdateRound(db *gorm.DB) {
	// several loops
	// first loop, update transaction status only depending on its chain status
	// we only need transaction block
	// TODO: batch get
	var transactions []models.ChainTransaction
	err := db.Find(&transactions, "is_stable != true").Error
	// err := db.Find(&transactions).Error
	if err != nil {
		logrus.WithError(err).Errorf("failed to get transactions")
		return
	}
	for _, tx := range transactions {
		if !tx.TxStatus.IsSent() {
			continue
		}
		fromAccount, err := tx.GetTransactionFrom(db)
		if err != nil {
			logrus.WithField("txId", tx.ID).WithError(err).Errorf("failed to get transaction from")
			return
		}
		meta, err := utils.W3.GetTransactionResult(fromAccount.Chain.Name, tx.Hash)
		if err != nil {
			logrus.WithField("txId", tx.ID).WithError(err).Errorf("failed to get transaction detail")
			return
		}
		// update Block number
		if meta.BlockNumber == nil {
			tx.BlockNumber = nil
		} else {
			blockNumber := meta.BlockNumber.ToInt().Uint64()
			tx.BlockNumber = &blockNumber
		}

		txNewStatus, err := types.InferSentTransactionStatus(meta, tx.TxStatus, fromAccount.Chain.GetTaggedBlockNumbers())
		if err != nil {
			logrus.WithField("txId", tx.ID).WithError(err).Errorf("failed to infer transaction status")
			return
		}
		if txNewStatus == tx.TxStatus {
			continue
		}
		// update status
		tx.TxStatus = txNewStatus
		// tx might be also stable if conflict nonce tx is finalized
		// if tx.TxStatus.IsStable() {
		// 	// update IsStable
		// 	tx.IsStable = true
		// }
		err = db.Save(&tx).Error
		if err != nil {
			logrus.WithField("txId", tx.ID).WithError(err).Errorf("failed to update transaction status")
			return
		}
		logrus.WithField("txId", tx.ID).WithField("service", "onchainstatus").Infof("transaction status updated to %d", tx.TxStatus)
	}
	// TODO: second loop, update error transaction status according to same nonce status to check if the transaction is already stable
}

func updateTransactionStatus(db *gorm.DB, tx *models.ChainTransaction) error {
	backupTx := *tx
	err := func() error {
		fromAccount, err := tx.GetTransactionFrom(db)
		if err != nil {
			return err
		}
		meta, err := utils.W3.GetTransactionResult(fromAccount.Chain.Name, tx.Hash)
		if err != nil {
			return errors.New("failed to get transaction detail")
		}
		// update Block number
		if meta.BlockNumber == nil {
			tx.BlockNumber = nil
		} else {
			blockNumber := meta.BlockNumber.ToInt().Uint64()
			tx.BlockNumber = &blockNumber
		}

		txNewStatus, err := types.InferSentTransactionStatus(meta, tx.TxStatus, fromAccount.Chain.GetTaggedBlockNumbers())
		if err != nil {
			return errors.New("failed to infer transaction status")
		}
		if txNewStatus == tx.TxStatus {
			return nil
		}
		// update status
		tx.TxStatus = txNewStatus
		// tx might be also stable if conflict nonce tx is finalized
		// if tx.TxStatus.IsStable() {
		// 	// update IsStable
		// 	tx.IsStable = true
		// }
		return nil
	}()
	if err != nil {
		*tx = backupTx
		return err
	}

	return nil
}

func StartTransactionChainStatusUpdateWorkerRound(ctx *QueueContext, maxBatchSize int) error {
	// several loops
	// first loop, update transaction status only depending on its chain status
	// we only need transaction block
	// TODO: batch get
	txs := ctx.PoolOrChainQueue.MustDequeBatch(maxBatchSize)
	if len(txs) == 0 {
		return nil
	}
	for i := len(txs) - 1; i >= 0; i-- {
		if !txs[i].TxStatus.IsSent() {
			panic("unexpected")
		}

		err := func() error {
			err := updateTransactionStatus(ctx.Db, txs[i])
			return err
		}()

		if err != nil {
			ctx.ErrQueue.MustEnqueWithLog(*txs[i], "onchainstatus", "error update transaction")
			// remove 1 transaction
			txs = append(txs[:i], txs[i+1:]...)
			continue
		}
	}

	// TODO: second loop, update error transaction status according to same nonce status to check if the transaction is already stable

	// save status
	err := models.SaveWithRetry(ctx.Db, &txs)
	if err != nil {
		logrus.WithError(err).Error("error saving transaction")
		// return err
		return err
	}
	for _, tx := range txs {
		if tx.TxStatus.IsStable() {
			logrus.WithField("txId", tx.ID).WithField("service", "onchainstatus").Info("transaction finalized")
			continue
		}
		ctx.PoolOrChainQueue.MustEnqueWithLog(*tx, "onchainstatus", "tx not finalized, watch another round")
	}
	return nil
}

// watch the pending transactions. move them into error states if certain circumstances met
// 1. current gas price & tx gas price & pending time（might vary depending on strategy） -> TxPoolLowPrice
// 2. ...
// do nothing here
func StartPendingStatusMoveRound(db *gorm.DB) {
	// skip now
}
