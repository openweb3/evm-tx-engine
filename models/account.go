package models

import (
	"time"

	"gorm.io/gorm"
)

// Account struct defines the fields of an account
type Account struct {
	gorm.Model
	// Public key of the account
	PublicKey string `gorm:"type:varchar(255)"`
	// Chains the account is associated with
	ChainAccounts []ChainAccount
	// Alias of the account
	Alias string `gorm:"type:varchar(255)"`
	// Status of the account locked,...
	AccountStatus string `gorm:"type:varchar(255)"`
	Impl          string `gorm:"type:varchar(255)"` // which implementation should use
}

// return ethaddress from public key
func (acct *Account) EthAddress() (string, error) {
	return "", nil
}

// 写竞争： Sponsor Update & Nonce Update. 但修改的不是同一字段，可能没有竞争？
// TODO：remove Sponsor Update （通过查询Sponsor Transaction来获取这两个值）
type ChainAccount struct {
	gorm.Model
	AccountId uint `gorm:"type:int"`
	// Account   Account `gorm:"foreignKey:AccountId"`
	// Address of the chain account
	Address             string `gorm:"type:varchar(255)"`
	ChainId             uint   `gorm:"type:int"` // the internal ChainId
	ChainType           string `gorm:"type:varchar(255)"`
	Chain               Chain  `gorm:"foreignKey:ChainId,ChainType;References:ID,Type"`
	LastSponsorInit     time.Time
	LastSponsorReceived time.Time
	// internal transaction count
	InternalNonce  uint64  `gorm:"type:int"`
	LatestNonce    *uint64 `gorm:"type:int"`
	SafeNonce      *uint64 `gorm:"type:int"`
	FinalizedNonce *uint64 `gorm:"type:int"`
}
