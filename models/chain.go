package models

import "github.com/openweb3/evm-tx-engine/utils"

// Chain struct defines the fields of a chain
type Chain struct {
	ID   uint   `gorm:"primaryKey"` // 先只考虑是primary key
	Type string `gorm:"type:varchar(255)"`
	// Name of the chain
	Name string `gorm:"type:varchar(255)"`
	// Boolean value to indicate if the chain is a testnet
	IsTestnet      bool   `gorm:"type:bool"`
	LatestBlock    uint64 `gorm:"type:int"`
	SafeBlock      uint64 `gorm:"type:int"`
	FinalizedBlock uint64 `gorm:"type:int"`
}

func (chain *Chain) GetTaggedBlockNumbers() utils.TaggedBlockNumbers {
	return utils.TaggedBlockNumbers{
		LatestBlock:    chain.LatestBlock,
		SafeBlock:      chain.SafeBlock,
		FinalizedBlock: chain.FinalizedBlock,
	}
}
