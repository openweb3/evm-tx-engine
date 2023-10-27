package models

import (
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	// Chain         string `gorm:"type:varchar(255)"`
	// From          string `gorm:"type:varchar(255)"`
	From          ChainAccount `gorm:"foreignKey:FromId"`
	FromId        uint         `gorm:"type:int"`
	TaskStatus    string       `gorm:"type:varchar(255)"`
	Fee           string       `gorm:"type:varchar(255)"`
	Hash          string       `gorm:"type:varchar(255)"`
	PreviousTask  int          `gorm:"type:varchar(255)"`
	PriceStrategy string       `gorm:"type:varchar(255)"`
	Priority      int          `gorm:"type:int"`

	RetryMaxAttempts int                `gorm:"type:int"`
	RetryDeadline    int64              `gorm:"type:bigint"`
	RetryInterval    int                `gorm:"type:int"`
	History          []ChainTransaction `gorm:"foreignKey:TaskId"` // Assuming HistoryItem is another struct

	FieldId int   `gorm:"type:uint"`
	Field   Field `gorm:"foreignKey:FieldId"`
	// FieldTo   string `gorm:"type:varchar(255)"`
	// FieldData string `gorm:"type:varchar(255)"`
}

func GetTask(db *gorm.DB, id uint) (*Task, error) {
	var task Task
	err := db.First(&task, id).Error
	return &task, err
}
