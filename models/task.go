package models

import (
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	// Chain         string `gorm:"type:varchar(255)"`
	// From          string `gorm:"type:varchar(255)"`
	From          ChainAccount
	TaskStatus    string `gorm:"type:varchar(255)"`
	Fee           string `gorm:"type:varchar(255)"`
	Hash          string `gorm:"type:varchar(255)"`
	PreviousTask  int    `gorm:"type:varchar(255)"`
	PriceStrategy string `gorm:"type:varchar(255)"`
	Priority      int    `gorm:"type:int"`

	RetryMaxAttempts int                `gorm:"type:int"`
	RetryDeadline    int64              `gorm:"type:bigint"`
	RetryInterval    int                `gorm:"type:int"`
	History          []ChainTransaction `gorm:"many2many:task_historys"` // Assuming HistoryItem is another struct

	FieldId int   `gorm:"type:int"`
	Field   Field `gorm:"foreignKey:FieldId"`
	// FieldTo   string `gorm:"type:varchar(255)"`
	// FieldData string `gorm:"type:varchar(255)"`
}
