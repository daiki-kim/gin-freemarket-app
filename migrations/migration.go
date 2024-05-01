package main

import (
	"gin-freemarket-app/infra"
	"gin-freemarket-app/models"
)

func main() {
	infra.Initialize()
	db := infra.SetupDB()

	if err := db.AutoMigrate(&models.Item{}, &models.User{}); err != nil {
		panic("failed to migrate database")
	}
}
