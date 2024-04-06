package main

import (
	"ginchat/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db, _ := gorm.Open(mysql.Open("root:123456@tcp(127.0.0.1:3306)/ginchat?charset=utf8mb4&parseTime=true&loc=Local"))
	db.AutoMigrate(&models.Community{})
	//utils.DB.AutoMigrate(&models.Message{})
}
