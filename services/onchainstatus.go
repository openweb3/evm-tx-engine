package services

import (
	"github.com/openweb3/evm-tx-engine/models"
	"github.com/openweb3/evm-tx-engine/utils"
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
			chain.LatestBlock = taggedBlockNumbers.LatestBlock
			chain.SafeBlock = taggedBlockNumbers.SafeBlock
			chain.FinalizedBlock = taggedBlockNumbers.FinalizedBlock
			db.Save(&chain)
		}
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
	err := db.Preload("Task.From.Chain").Where("is_stable == true").Find(&transactions).Error
	if err != nil {
		return
	}
	for _, tx := range transactions {
		if !tx.TxStatus.IsSent() {
			continue
		}
		meta, err := utils.W3.GetTransactionDetail(tx.Task.From.Chain.Name, tx.Hash)
		if err == nil {
			return
		}
		// update Block number
		tx.BlockNumber = uint((*meta.BlockNumber).ToInt().Int64())
		txNewStatus, err := utils.InferSentTransactionStatus(meta, tx.TxStatus, tx.Task.From.Chain.GetTaggedBlockNumbers())
		if err == nil {
			return
		}
		// update status
		tx.TxStatus = txNewStatus
		// tx might be also stable if conflict nonce tx is finalized
		if tx.TxStatus.IsStable() {
			// update IsStable
			tx.IsStable = true
		}
		db.Save(&tx)
	}
	// TODO: second loop, update error transaction status according to same nonce status to check if the transaction is already stable
}

// // check all unstable transaction status, including PoolError
// // how to check
// func StartChainStatusUpdater(db *gorm.DB) {

// 	unstableChainStatus := []utils.TxStatus{utils.TxPoolPending, utils.TxPoolError, utils.TxChainLatest, utils.TxChainSafe, utils.TxChainLatestError, utils.TxChainSafeError}

// 	var unstableChainTxs []models.ChainTransaction

// 	// TODO: consider alternative implementation
// 	//       query each transaction of from with difference nonces

// 	// filterOut unstable transactions
// 	db.Model(&models.ChainTransaction{}).Where("status IN ?", unstableChainStatus).Find(&unstableChainTxs)

// 	// TODO: batch query & update
// 	for _, tx := range unstableChainTxs {
// 		tx.UpdateStatus(db)
// 	}
// }

// watch the pending transactions. move them into error states if certain circumstances met
// 1. current gas price & tx gas price & pending time（might vary depending on strategy） -> TxPoolLowPrice
// 2. ...
// do nothing here
func StartPendingStatusMoveRound(db *gorm.DB) {
	// skip now
}
