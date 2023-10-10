package services

import (
	"github.com/openweb3/evm-tx-engine/models"
	"gorm.io/gorm"
)

// pick task and add tx transactions
// the task should have no related tx
func StartPickerService(db *gorm.DB) {
	var tasks []models.Task
	pageSize := 100

	// 使用 gorm 从数据库中选择满足以下条件的任务：

	// 任务没有相关的交易。
	// 按 From 分组，并在每个组中选择 Priority 最高的任务。
	// 在所有选中的任务中，选择 CreatedAt 最早的任务。
	// 每次选择一个完整的分页（例如100个任务）。
	// 为每个选中的任务添加一个新的 ChainTransaction，并设置以下字段：

	// TaskId 为任务的ID。
	// IsCancelTransaction 为 false。
	// TxStatus 为 "TARGET_QUEUE"。
	// 其他字段为空。

	// TODO: lock
	// TODO: modify task status
	db.Order("From, Priority, CreatedAt").Where("NOT EXISTS (SELECT 1 FROM chain_transactions WHERE chain_transactions.TaskId = tasks.ID)").Limit(pageSize).Find(&tasks)

	for _, task := range tasks {
		// 为任务创建一个新的 ChainTransaction
		tx := models.ChainTransaction{
			TaskId:              task.ID,
			IsCancelTransaction: false,
			TxStatus:            "TARGET_QUEUE",
		}

		// 将交易保存到数据库
		db.Create(&tx)
	}
}
