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
	ChainAccount []ChainAccount
	// Alias of the account
	Alias string `gorm:"type:varchar(255)"`
	// Status of the account
	AccountStatus string `gorm:"type:varchar(255)"`
	Impl          string `gorm:"type:varchar(255)"` // which implementation should use
}

// // AccountChain serves as a join table for the many-to-many relationship between Account and Chain
// // NOTE: account_chain should be unique
// type AccountChain struct {
// 	AccountID           uint
// 	ChainID             uint // the internal ChainID
// 	LastSponsorInit     time.Time
// 	LastSponsorReceived time.Time
// 	// internal transaction count
// 	InternalNonce  uint `gorm:"type:int"`
// 	LatestNonce    uint `gorm:"type:int"`
// 	SafeNonce      uint `gorm:"type:int"`
// 	FinalizedNonce uint `gorm:"type:int"`
// }

// 写竞争： Sponsor Update & Nonce Update. 但修改的不是同一字段，可能没有竞争？
// TODO：remove Sponsor Update （通过查询Sponsor Transaction来获取这两个值）
type ChainAccount struct {
	gorm.Model
	AccountID           uint
	ChainID             uint // the internal ChainID
	Chain               Chain
	LastSponsorInit     time.Time
	LastSponsorReceived time.Time
	// internal transaction count
	InternalNonce  uint `gorm:"type:int"`
	LatestNonce    uint `gorm:"type:int"`
	SafeNonce      uint `gorm:"type:int"`
	FinalizedNonce uint `gorm:"type:int"`
}