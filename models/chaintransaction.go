package models

import (
	"github.com/openweb3/evm-tx-engine/utils"
	"gorm.io/gorm"
)

// TODO: field type
type ChainTransaction struct {
	gorm.Model
	TaskId              uint           `gorm:"type:int"`
	Task                Task           `gorm:"foreignKey:TaskId"`
	IsCancelTransaction bool           `gorm:"type:bool"` // if this is the cancel transaction, then it will not use the task id attached fields to construct the transaction
	IsStable            bool           `gorm:"type:bool"`
	BlockNumber         uint           `gorm:"type:int"`
	Raw                 []byte         `gorm:"type:TINYBLOB"`     // the raw transaction
	Hash                string         `gorm:"type:VARCHAR(64)"`  // Transaction Hash
	TxStatus            utils.TxStatus `gorm:"type:varchar(255)"` // Transaction Stage
	ErrorMessage        string         `gorm:"type:varchar(255)"` // Error of the transaction(if met)
	ErrorCode           string         `gorm:"type:varchar(255)"`
	FieldId             uint           `gorm:"type:int"`
	Field               Field          `gorm:"foreignKey:FieldId"`
}

func FetchChainTransactionStatusAndStabilityFromChain(db *gorm.DB, chainTransaction *ChainTransaction) (utils.TxStatus, bool, error) {
	// read chain from chainTransaction
	// chain
	return utils.TxChainSafe, false, nil

}

// TODO: prepare the transaction for signer to sign
func (chainTransaction *ChainTransaction) PrepareTransactionToSign() (interface{}, error) {
	return nil, nil
}

// TODO: check if could sign for specific chain?
func (chainTransaction *ChainTransaction) GetAccount() (*Account, error) {
	return nil, nil
}

// func GetNonceStability()
