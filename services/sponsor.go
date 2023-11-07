package services

// temporary implementation
// BalanceCheck will let transaction pass if balance is enough
// 目前是个空实现

func StartBalanceCheckWorkerRound(ctx *QueueContext, maxBatchSize int) error {
	txs := ctx.BalanceCheckingQueue.MustDequeBatch(maxBatchSize)

	for _, tx := range *txs {
		// checks balance enough
		SetGasPrice(&tx)
		ctx.NonceManagingQueue.MustEnqueWithLog(tx, "BalanceCheck", "moved to GasEnough Queue")
	}
	return nil
}

// sponsor 错误处理
// func StartSponsorService(db *gorm.DB) {
// 	errorCodes := []utils.TxErrorCode{utils.TxInternalInsufficientBalance, utils.TxPoolInsufficientBalance}

// 	// 根据提供的错误码从数据库中获取唯一的 'From' 地址
// 	var addresses []string
// 	db.Model(&models.ChainTransaction{}).Where("ErrorCode IN ?", errorCodes).Distinct("From").Pluck("From", &addresses)

// 	for _, address := range addresses {
// 		// 获取给定 'From' 地址的最新的 UpdatedAt
// 		var latestTx models.ChainTransaction
// 		db.Where("From = ? AND ErrorCode IN ?", address, errorCodes).Order("UpdatedAt DESC").First(&latestTx)

// 		// 检查 'From' 地址的 LastSponsorReceived
// 		var chainAccount models.ChainAccount
// 		db.Where("Address = ? AND ChainID = ?", address, chainTransaction.Task.Chain.ID).First(&account)

// 		// 检查 ChainTransaction 的最新的 UpdatedAt 是否在 LastSponsorReceived 之后
// 		if latestTx.UpdatedAt.After(account.LastSponsorReceived) && account.LastSponsorReceived.After(account.LastSponsorInit) {
// 			// TODO: 向地址发送交易费
// 			// 这只是一个占位符，发送费用的实际逻辑需要实现

// 			// 更新 'From' 地址的 LastSponsorInit 为当前时间戳
// 			account.LastSponsorInit = time.Now()
// 			db.Save(&account)
// 		}
// 	}
// }
