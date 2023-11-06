package services

import (
	"fmt"

	"github.com/openweb3/evm-tx-engine/models"
	"github.com/openweb3/evm-tx-engine/utils"
	"github.com/sirupsen/logrus"
)

// pick task and add tx transactions
// the task should have no related tx
func StartPickerRound(ctx *QueueContext, maxSize int) error {
	var tasks []models.Task
	pageSize := maxSize

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
	// db.Order("From, Priority, CreatedAt").Where("NOT EXISTS (SELECT 1 FROM chain_transactions WHERE chain_transactions.TaskId = tasks.ID)").Preload("Field").Limit(pageSize).Find(&tasks)

	// =================================
	// First of all use a simple implementation here
	// search all tasks with no chain_transaction and then add a transaction for the task
	ctx.Db.Where("NOT EXISTS (SELECT 1 FROM chain_transactions WHERE chain_transactions.task_id = tasks.ID)").Preload("Field").Limit(pageSize).Find(&tasks)

	for _, task := range tasks {
		// 为任务创建一个新的 ChainTransaction
		// field :=
		// db.Create(&field)

		tx := models.ChainTransaction{
			TaskId:              task.ID,
			IsCancelTransaction: false,
			TxStatus:            utils.TxInternalTargetQueue,
			Field: models.Field{
				To:                task.Field.To,
				MaxFeePerGas:      task.Field.MaxFeePerGas,
				Data:              task.Field.Data,
				GasLimit:          task.Field.GasLimit,
				GasPrice:          task.Field.GasPrice,
				Value:             task.Field.Value,
				PriorityFeePerGas: task.Field.PriorityFeePerGas,
			},
		}

		// 将交易保存到数据库
		err := ctx.Db.Save(&tx).Error
		if err != nil {
			logrus.WithError(err).Error("Failed to create tx")
			return err
		}
		ctx.TargetQueue.MustEnqueWithLog(tx, "Picker", fmt.Sprintf("transaction created for task %d", task.ID))
	}
	return nil
}
