package models

import "gorm.io/gorm"

// TODO: field type
type ChainTransaction struct {
	gorm.Model
	TaskId       int    `gorm:"type:int"`          // TODO: foreign key
	Raw          []byte `gorm:"type:TINYBLOB"`     // the raw transaction
	Hash         string `gorm:"type:VARCHAR(64)"`  // Transaction Hash
	TxStatus     string `gorm:"type:varchar(255)"` // Transaction Stage
	ErrorMessage string `gorm:"type:varchar(255)"` // Error of the transaction(if met)
	ErrorCode    string `gorm:"type:varchar(255)"`
}
