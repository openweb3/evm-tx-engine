package models

import "github.com/holiman/uint256"

type Field struct {
	To string `gorm:"type:varchar(255)"`
	// Function      string             `gorm:"type:varchar(255)"`
	// Params        []string           `gorm:"type:json"`
	MaxFeePerGas      *uint256.Int `gorm:"default:0;not null"`
	Data              []byte       `gorm:"type:TINYBLOB"`
	Nonce             *uint64      `gorm:"type:int"`
	GasLimit          *uint256.Int `gorm:"default:0;not null"`
	GasPrice          *uint256.Int `gorm:"default:0;not null"`
	Value             *uint256.Int `gorm:"default:0;not null"`
	PriorityFeePerGas *uint256.Int `gorm:"default:0;not null"`
}
