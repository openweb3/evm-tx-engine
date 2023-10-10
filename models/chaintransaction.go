package models

import "gorm.io/gorm"

// TODO: field type
type ChainTransaction struct {
	gorm.Model
	TaskId              uint   `gorm:"type:int"`  // TODO: foreign key
	IsCancelTransaction bool   `gorm:"type:bool"` // if this is the cancel transaction, then it will not use the task id attached fields to construct the transaction
	IsStable            bool   `gorm:"type:bool"`
	Raw                 []byte `gorm:"type:TINYBLOB"`     // the raw transaction
	Hash                string `gorm:"type:VARCHAR(64)"`  // Transaction Hash
	TxStatus            string `gorm:"type:varchar(255)"` // Transaction Stage
	ErrorMessage        string `gorm:"type:varchar(255)"` // Error of the transaction(if met)
	ErrorCode           string `gorm:"type:varchar(255)"`
	FieldId             uint    `gorm:"type:int"`
	Field               Field
}
