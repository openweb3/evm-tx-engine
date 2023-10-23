package utils

import (
	"time"

	"gorm.io/gorm"
)

type serviceRound func(db *gorm.DB)

func StartService(db *gorm.DB, intervalMs uint, roundFunc serviceRound) {
	for {
		roundFunc(db)
		time.Sleep(time.Duration(intervalMs) * time.Millisecond)
	}
}
