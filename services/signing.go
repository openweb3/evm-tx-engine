package services

import (
	"github.com/openweb3/evm-tx-engine/accountadapter"
	"github.com/openweb3/evm-tx-engine/models"
	"github.com/openweb3/evm-tx-engine/utils"
	"gorm.io/gorm"
)

// pop a batch of transactions from constructed, and send them to the signing service. Modify the `Raw` field and modify the status field
func StartSigningRound(db *gorm.DB) {
	var txs []models.ChainTransaction
	db.Where("tx_status = ?", utils.TxInternalConstructed).Find(&txs)
	// will modify txs
	err := batchSign(&txs)
	if err == nil {
		// do something
		return
	}
	for _, tx := range txs {
		tx.TxStatus = utils.TxInternalSigned
	}
	db.Save(&txs)
}

func batchSign(txs *[]models.ChainTransaction) error {
	for _, tx := range *txs {
		txReadyToSign, err := tx.PrepareTransactionToSign()
		if err != nil {
			return err
		}
		account, err := tx.GetAccount()
		adapter, err := accountadapter.GetAccountAdapter(account)
		signed, err := adapter.SignTransaction(account.Address, txReadyToSign)
		tx.Raw = signed
	}
	return nil
}
