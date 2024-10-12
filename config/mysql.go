package config

import (
	"gorm.io/gorm"
	"wallet/pkg/db/dbconn"
	"wallet/pkg/zlogger"
)

var DB *gorm.DB

func InitDB() {
	db, err := dbconn.NewGormDB()
	if err != nil {
		zlogger.Errorf("Error connecting to database: %s", err.Error())
		panic(err)
	}

	DB = db
	zlogger.Info("Connected to database successfully")
}
