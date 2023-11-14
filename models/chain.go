package models

import (
	"github.com/openweb3/evm-tx-engine/types"
	"gorm.io/gorm"
)

// Chain struct defines the fields of a chain
type Chain struct {
	gorm.Model
	ChainId uint   `gorm:"index:id_type_idx,unique"` // 先只考虑是primary key
	Type    string `gorm:"index:id_type_idx,unique;type:varchar(255)"`
	// Name of the chain
	Name string `gorm:"type:varchar(255)"`
	// Boolean value to indicate if the chain is a testnet
	IsTestnet      bool   `gorm:"type:bool"`
	LatestBlock    uint64 `gorm:"type:int"`
	SafeBlock      uint64 `gorm:"type:int"`
	FinalizedBlock uint64 `gorm:"type:int"`
}

func (chain *Chain) GetTaggedBlockNumbers() types.TaggedBlockNumbers {
	return types.TaggedBlockNumbers{
		LatestBlock:    chain.LatestBlock,
		SafeBlock:      chain.SafeBlock,
		FinalizedBlock: chain.FinalizedBlock,
	}
}
