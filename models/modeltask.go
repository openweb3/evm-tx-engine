package models

import (
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Chain         string `gorm:"type:varchar(255)"`
	From          string `gorm:"type:varchar(255)"`
	TaskStatus    string `gorm:"type:varchar(255)"`
	Fee           string `gorm:"type:varchar(255)"`
	Hash          string `gorm:"type:varchar(255)"`
	PreviousTask  int    `gorm:"type:varchar(255)"`
	PriceStrategy string `gorm:"type:varchar(255)"`
	Priority      int    `gorm:"type:int"`

	RetryMaxAttempts int                `gorm:"type:int"`
	RetryDeadline    int64              `gorm:"type:bigint"`
	RetryInterval    int                `gorm:"type:int"`
	History          []ChainTransaction `gorm:"type:json"` // Assuming HistoryItem is another struct

	FieldId int `gorm:"type:int"`
	Field   Field
	// FieldTo   string `gorm:"type:varchar(255)"`
	// FieldData string `gorm:"type:varchar(255)"`
}

type Field struct {
	To string `gorm:"type:varchar(255)"`
	// Function      string             `gorm:"type:varchar(255)"`
	// Params        []string           `gorm:"type:json"`
	MaxFeePerGas int    `gorm:"type:int"`
	Data         []byte `gorm:"type:TINYBLOB"`
}
