package models

import (
	"time"

	"gorm.io/gorm"
)

// Account struct defines the fields of an account
type Account struct {
	gorm.Model
	// Address of the account
	Address string `gorm:"type:varchar(255)"`
	// Public key of the account
	PublicKey string `gorm:"type:varchar(255)"`
	// Chains the account is associated with
	Chains []Chain `gorm:"many2many:account_chains"`
	// Alias of the account
	Alias string `gorm:"type:varchar(255)"`
	// Status of the account
	AccountStatus string `gorm:"type:varchar(255)"`
	// internal transaction count
	TransactionCountInternal uint `gorm:"type:int"`
	LastSponsorInit          time.Time
	LastSponsorReceived      time.Time
}

// Chain struct defines the fields of a chain
type Chain struct {
	ID   uint   `gorm:"primaryKey"` // 先只考虑是primary key
	Type string `gorm:"type:varchar(255)"`
	// Name of the chain
	Name string `gorm:"type:varchar(255)"`
	// Boolean value to indicate if the chain is a testnet
	IsTestnet bool `gorm:"type:bool"`
}

// AccountChain serves as a join table for the many-to-many relationship between Account and Chain
type AccountChain struct {
	AccountID uint
	ChainID   uint
}
