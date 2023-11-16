package models

import (
	"github.com/openweb3/evm-tx-engine/types/code"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// TODO: field type
type ChainTransaction struct {
	gorm.Model
	TaskId              uint `gorm:"type:int"`
	IsCancelTransaction bool `gorm:"type:bool"` // if this is the cancel transaction, then it will not use the task id attached fields to construct the transaction
	// IsStable            bool          `gorm:"type:bool"`
	BlockNumber  *uint64       `gorm:"type:int"`
	Raw          *[]byte       `gorm:"type:TINYBLOB"`     // the raw transaction
	Hash         string        `gorm:"type:VARCHAR(64)"`  // Transaction Hash
	TxStatus     code.TxStatus `gorm:"type:varchar(255)"` // Transaction Stage
	ErrorMessage string        `gorm:"type:varchar(255)"` // Error of the transaction(if met)
	ErrorCode    string        `gorm:"type:varchar(255)"`
	Field
}

func FetchChainTransactionStatusAndStabilityFromChain(db *gorm.DB, chainTransaction *ChainTransaction) (code.TxStatus, bool, error) {
	// read chain from chainTransaction
	// chain
	return code.TxChainSafe, false, nil

}

// TODO: check if could sign for specific chain?
func (chainTransaction *ChainTransaction) GetTransactionFrom(db *gorm.DB) (*ChainAccount, error) {
	var task Task
	err := db.Preload("From.Chain").Model(&Task{}).Where("id = ?", chainTransaction.TaskId).First(&task).Error
	return &task.From, err
}

func (chainTransaction *ChainTransaction) GetSigner(db *gorm.DB) (*Account, error) {
	var account Account
	chainAccount, err := chainTransaction.GetTransactionFrom(db)
	if err != nil {
		return nil, err
	}
	err = db.Preload(clause.Associations).First(&account, chainAccount.AccountId).Error
	return &account, err
}

// func GetNonceStability()
