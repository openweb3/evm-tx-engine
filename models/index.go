package models

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/openweb3/evm-tx-engine/config"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDB(dbName string) *gorm.DB {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	var err error
	config.Init()
	dbConfig := config.GetConfig().Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Migrate the schema
	err = Migrate(db)

	if err != nil {
		panic(err)
	}
	return db
}

func Migrate(db *gorm.DB) error {
	// NOTE: DON'T TRY TO MERGE THE 2 AUDOMIGRATE
	// merging them will bring about errors
	err := db.AutoMigrate(
		&Chain{},
	)
	if err != nil {
		return err
	}
	err = db.AutoMigrate(
		// &Chain{},
		&Account{},
		&ChainAccount{},
		&Task{},
		&ChainTransaction{},
	)
	return err
}

func SaveWithRetry(db *gorm.DB, value interface{}) error {
	defaultRetryInterval := 1000
	maxRetry := 3
	for i := 0; i < int(maxRetry); i++ {
		err := db.Save(value).Error
		if err == nil {
			return nil
		}
		time.Sleep(time.Duration(defaultRetryInterval) * time.Millisecond)
	}
	err := fmt.Errorf("data failed to save after %d retries: %+v", maxRetry, value)
	logrus.WithError(err).WithField("stack", string(debug.Stack())).Errorf("failed to save data")
	return err
}
