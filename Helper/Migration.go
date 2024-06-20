package Helper

import (
	"ims/Database"
	"ims/Model"
	"log"
)

func AutoMigrate() {
	Database.Connect()
	Database.Database.AutoMigrate(&Model.User{})
	Database.Database.AutoMigrate(&Model.Product{})
	Database.Database.Exec("CREATE TYPE changetype AS ENUM ('addition', 'subtraction', 'insert', 'delete', 'update');")
	Database.Database.AutoMigrate(&Model.Activity{})
	log.Println("Migration Success")
}
