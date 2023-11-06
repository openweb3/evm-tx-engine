package services

import (
	"errors"

	"github.com/enriquebris/goconcurrentqueue"
	"github.com/openweb3/evm-tx-engine/models"
	"github.com/sirupsen/logrus"
)

// element should be type of models.ChainTransaction
// do not use pointer to prevent element modification after enque
type ChainTransactionQueue goconcurrentqueue.FIFO

func NewChainTransactionQueue() *ChainTransactionQueue {
	return (*ChainTransactionQueue)(goconcurrentqueue.NewFIFO())
}

// returns error if no elements can be deque
func (queue *ChainTransactionQueue) Deque() (models.ChainTransaction, error) {
	tx, err := (*goconcurrentqueue.FIFO)(queue).Dequeue()
	// returns error if empty or locked
	// then we will return a nil
	if err != nil {
		return models.ChainTransaction{}, err
	}
	tx_, ok := tx.(models.ChainTransaction)
	if !ok {
		// should never happen
		panic("element type unexpected")
	}
	// check again to ensure data validity
	// could remove if there would be urgent performance requirements
	if tx_.Field.ID == 0 {
		// should never happen
		panic("element's Field is not attached")
	}
	return tx_, nil
}

func (queue *ChainTransactionQueue) Enque(tx models.ChainTransaction) error {
	// transaction field should be attached to the transactions
	if tx.Field.ID == 0 {
		return errors.New("element's Field is not attached")
	}
	if tx.ID == 0 {
		return errors.New("transaction should have an attached id")
	}
	err := (*goconcurrentqueue.FIFO)(queue).Enqueue(tx)
	// Should not return error here because no plan to lock
	if err != nil {
		panic("unexpected lock happened")
	}
	return nil
}

// DequeBatch returns a slice of transactions from the queue, with maxSize limit
// might return empty array with 0 elements
func (queue *ChainTransactionQueue) MustDequeBatch(maxSize int) *[]models.ChainTransaction {
	var transactions []models.ChainTransaction = make([]models.ChainTransaction, 0, maxSize)

	for i := 0; i < maxSize; i++ {

		if tx, err := queue.Deque(); err == nil {
			transactions = append(transactions, tx)
		} else {
			// err means queue is already empty
			break
		}
	}
	return &transactions
}

// This function won't operate db
// Enque returning an error means the data is not valid
func (destQueue *ChainTransactionQueue) MustEnqueWithLog(tx models.ChainTransaction, workerName, successLog string) {
	err := destQueue.Enque(tx)
	if err != nil {
		logrus.WithField("txId", tx.ID).WithField("fieldId", tx.Field.ID).WithField("worker", workerName).WithError(err).Error("failed to enqueue transaction")
		return
	}
	logrus.WithField("txId", tx.ID).WithField("service", workerName).Info(successLog)
}
