package services

import (
	"time"

	"gorm.io/gorm"
)

type workerRound func(ctx *QueueContext, maxSize int) error

func startWorker(ctx *QueueContext, intervalMs uint, maxSize int, workerRoundFunc workerRound) {
	for {
		workerRoundFunc(ctx, maxSize)
		time.Sleep(time.Duration(intervalMs) * time.Millisecond)
	}
}

type serviceRound func(db *gorm.DB)

func startService(db *gorm.DB, intervalMs uint, roundFunc serviceRound) {
	for {
		roundFunc(db)
		time.Sleep(time.Duration(intervalMs) * time.Millisecond)
	}
}

func StartWorkers(ctx *QueueContext) {
	defaultInterval := 200
	maxSize := 128
	go startService(ctx.Db, uint(defaultInterval), StartTaggedBlockNumberUpdateRound)

	go startWorker(ctx, uint(defaultInterval), maxSize, StartPickerWorkerRound)
	go startWorker(ctx, uint(defaultInterval), maxSize, StartExecutionSimulationWorkerRound)
	go startWorker(ctx, uint(defaultInterval), maxSize, StartBalanceCheckWorkerRound)
	go startWorker(ctx, uint(defaultInterval), maxSize, StartNonceManageWorkerRound)
	go startWorker(ctx, uint(defaultInterval), maxSize, StartSigningWorkerRound)
	go startWorker(ctx, uint(defaultInterval), maxSize, StartSenderWorkerRound)
	go startWorker(ctx, uint(defaultInterval), maxSize, StartTransactionChainStatusUpdateWorkerRound)

}

// TODO: init queues from database rather than empty queues
// func InitQueues() {
// 	queue := goconcurrentqueue.NewFIFO()
// }
