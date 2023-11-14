package services

func StartExecutionSimulationWorkerRound(ctx *QueueContext, maxBatchSize int) error {
	txs := ctx.TargetQueue.MustDequeBatch(maxBatchSize)

	for _, tx := range txs {

		if tx.Field.GasLimit.CmpUint64(0) != 0 {
			continue
		}
		tx.Field.GasLimit.SetUint64(10000000)
		ctx.BalanceCheckingQueue.MustEnqueWithLog(*tx, "ExecutionSimulator", "gas limit set")
	}
	return nil
}
