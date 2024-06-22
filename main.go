package main

import (
	"ims/Database"
	helper "ims/Helper"
	"ims/Router"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	loadEnv()
	Database.Connect()

	// Migration Database
	helper.AutoMigrate()

	Router.SetupRouter().Run(":8080")
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
