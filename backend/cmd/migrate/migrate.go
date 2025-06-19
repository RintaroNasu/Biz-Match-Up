package main

import (
	"backend/internal/db"
	"backend/internal/domain/model"

	"fmt"
)

func main() {
	dbConn := db.NewDB()
	defer fmt.Println("Successfully Migrated")
	defer db.CloseDB(dbConn)
	dbConn.AutoMigrate(&model.User{}, &model.Reason{})
}
