package services_test

import (
	"testing"

	"github.com/openweb3/evm-tx-engine/models"
	"github.com/openweb3/evm-tx-engine/services"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestNewChainTransactionQueue(t *testing.T) {
	queue := services.NewChainTransactionQueue()
	assert.NotNil(t, queue)
	// Add additional checks if needed, e.g., checking initial size, etc.
}

func TestDeque(t *testing.T) {
	queue := services.NewChainTransactionQueue()

	// Test dequeue from empty queue
	_, err := queue.Deque()
	assert.Error(t, err)

	// field := models.Field{
	// 	// gorm.Model{
	// 	ID: 1,
	// 	// },
	// }

	// Test successful dequeue
	tx := models.ChainTransaction{
		Model: gorm.Model{
			ID: 1,
		},
		Field: models.Field{
			Model: gorm.Model{
				ID: 1,
			}},
	} // Adjust this based on your model
	err = queue.Enque(tx)
	assert.NoError(t, err)
	dequeuedTx, err := queue.Deque()
	assert.NoError(t, err)
	assert.Equal(t, tx, dequeuedTx)
}

func TestEnque(t *testing.T) {
	queue := services.NewChainTransactionQueue()

	// Test enqueue with invalid element
	invalidTx := models.ChainTransaction{}
	err := queue.Enque(invalidTx)
	assert.Error(t, err)

	// Test successful enqueue
	validTx := models.ChainTransaction{
		Model: gorm.Model{
			ID: 1,
		},
		Field: models.Field{
			Model: gorm.Model{
				ID: 1,
			}},
	} // Adjust this based on your model
	err = queue.Enque(validTx)
	assert.NoError(t, err)
}

func TestMustDequeBatch(t *testing.T) {
	queue := services.NewChainTransactionQueue()

	// Test dequeue batch from empty queue
	txs := queue.MustDequeBatch(5)
	assert.Equal(t, 0, len(*txs))

	// Enqueue some transactions
	for i := 1; i <= 10; i++ {
		tx := models.ChainTransaction{
			Model: gorm.Model{
				ID: uint(i),
			},
			Field: models.Field{
				Model: gorm.Model{
					ID: uint(i),
				}},
		}
		queue.Enque(tx)
	}

	// Test dequeue batch with maxSize less than available elements
	txs = queue.MustDequeBatch(5)
	assert.Equal(t, 5, len(*txs))

	// Test dequeue batch with maxSize more than available elements
	txs = queue.MustDequeBatch(10)
	assert.Equal(t, 5, len(*txs)) // Only 5 elements should remain after the previous deque operation
}
