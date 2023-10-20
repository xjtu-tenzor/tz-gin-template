package model

import (
	"fmt"
	"log"
	"time"

	"template/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&&parseTime=True&loc=Local",
		config.Config.MysqlUser,
		config.Config.MysqlPass,
		config.Config.MysqlHost,
		config.Config.MysqlPort,
		config.Config.MysqlName)
	var dbLogger logger.Interface
	if config.DatabaseLogger == nil {
		dbLogger = logger.Default.LogMode(logger.Info)
	} else {
		dbLogger = logger.New(
			log.New(config.DatabaseLogger, "\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             200 * time.Millisecond,
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: true,
				Colorful:                  false,
			},
		)
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: dbLogger})
	if err != nil {
		panic(err)
	}

	DB = db

	if !config.Config.AppProd {
		initModel()
	}

}

func initModel() {

	// example
	// begin
	DB.AutoMigrate(&Resource{})
	//end

}
