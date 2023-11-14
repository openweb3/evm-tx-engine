package services

import (
	"errors"

	"gorm.io/gorm"
)

type QueueContext struct {
	Db                   *gorm.DB
	TargetQueue          *ChainTransactionQueue
	BalanceCheckingQueue *ChainTransactionQueue
	NonceManagingQueue   *ChainTransactionQueue
	SigningQueue         *ChainTransactionQueue
	SendingQueue         *ChainTransactionQueue
	PoolOrChainQueue     *ChainTransactionQueue
	ErrQueue             *ChainTransactionQueue
}

// *singleton
var queueContext *QueueContext

var queueInited bool

func GetQueueContext() (*QueueContext, error) {
	if queueInited {
		return queueContext, nil
	}
	return nil, errors.New("queue not initialized")
}

func InitQueueContext(db *gorm.DB) (*QueueContext, error) {
	if queueInited {
		return nil, errors.New("queue already initialized")
	}
	// rename
	queueContext = &QueueContext{
		Db:                   db,
		TargetQueue:          NewChainTransactionQueue(),
		BalanceCheckingQueue: NewChainTransactionQueue(),
		NonceManagingQueue:   NewChainTransactionQueue(),
		SigningQueue:         NewChainTransactionQueue(),
		SendingQueue:         NewChainTransactionQueue(),
		PoolOrChainQueue:     NewChainTransactionQueue(),
		ErrQueue:             NewChainTransactionQueue(),
	}
	queueInited = true
	return queueContext, nil
}
