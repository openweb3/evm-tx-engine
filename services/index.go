package services

import (
	"time"

	"github.com/openweb3/evm-tx-engine/utils"
)

// func StartServices(db *gorm.DB) {
// 	// 500ms as default round interval
// 	defaultInterval := 100
// 	go utils.StartService(db, uint(defaultInterval), StartPickerRound)

// 	go utils.StartService(db, uint(defaultInterval), StartExecutionSimulationRound)
// 	go utils.StartService(db, uint(defaultInterval), StartBalanceCheckRound)

// 	go utils.StartService(db, uint(defaultInterval), StartNonceManageRound)
// 	go utils.StartService(db, uint(defaultInterval), StartPriceManageRound)

// 	go utils.StartService(db, uint(defaultInterval), StartSigningRound)
// 	go utils.StartService(db, uint(defaultInterval), StartSenderRound)

// 	go utils.StartService(db, uint(defaultInterval), StartTaggedBlockNumberUpdateRound)
// 	go utils.StartService(db, uint(defaultInterval), StartTransactionChainStatusUpdateRound)
// 	go utils.StartService(db, uint(defaultInterval), StartPendingStatusMoveRound)

// }

// func StartServices(db *gorm.DB) {
// 	// 500ms as default round interval
// 	defaultInterval := 100
// 	go utils.StartService(db, uint(defaultInterval), StartPickerRound)

// 	go utils.StartService(db, uint(defaultInterval), StartExecutionSimulationRound)
// 	go utils.StartService(db, uint(defaultInterval), StartBalanceCheckRound)

// 	go utils.StartService(db, uint(defaultInterval), StartNonceManageRound)
// 	go utils.StartService(db, uint(defaultInterval), StartPriceManageRound)

// 	go utils.StartService(db, uint(defaultInterval), StartSigningRound)
// 	go utils.StartService(db, uint(defaultInterval), StartSenderRound)

// 	go utils.StartService(db, uint(defaultInterval), StartTaggedBlockNumberUpdateRound)
// 	go utils.StartService(db, uint(defaultInterval), StartTransactionChainStatusUpdateRound)
// 	go utils.StartService(db, uint(defaultInterval), StartPendingStatusMoveRound)

// }

type workerRound func(ctx *QueueContext, maxSize int) error

func startWorker(ctx *QueueContext, intervalMs uint, maxSize int, workerRoundFunc workerRound) {
	for {
		workerRoundFunc(ctx, maxSize)
		time.Sleep(time.Duration(intervalMs) * time.Millisecond)
	}
}

func StartWorkers(ctx *QueueContext) {
	defaultInterval := 200
	maxSize := 16
	go utils.StartService(ctx.Db, uint(defaultInterval), StartTaggedBlockNumberUpdateRound)

	go startWorker(ctx, uint(defaultInterval), maxSize, StartPickerRound)
	go startWorker(ctx, uint(defaultInterval), maxSize, StartExecutionSimulationWorkerRound)
	go startWorker(ctx, uint(defaultInterval), maxSize, StartBalanceCheckWorkerRound)
	go startWorker(ctx, uint(defaultInterval), maxSize, StartNonceManageWorkderRound)
	go startWorker(ctx, uint(defaultInterval), maxSize, StartSigningWorkerRound)
	go startWorker(ctx, uint(defaultInterval), maxSize, StartSenderWorkerRound)
	go startWorker(ctx, uint(defaultInterval), maxSize, StartTransactionChainStatusUpdateWorkerRound)

}

// func InitQueues() {
// 	queue := goconcurrentqueue.NewFIFO()
// }
