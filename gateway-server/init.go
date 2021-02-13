package main

import (
	"github.com/codenoid/pc-auth-notifier/shared-packages/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func init() {
	db, err := gorm.Open(sqlite.Open("pc-auth-notifier.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	mainDB = db

	// Migrate the schema
	db.AutoMigrate(&model.AuthLog{})
}
