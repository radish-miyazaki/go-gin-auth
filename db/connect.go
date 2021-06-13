package db

import (
	"github.com/radish-miyazaki/go-auth/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	// DBとの接続
	connection, err := gorm.Open(mysql.Open("root:password@/go_auth"), &gorm.Config{})
	if err != nil {
		panic("couldn't connect to the database!")
	}

	// TODO: gormのマイグレーションを他のライブラリに変更する
	DB = connection
	if err := connection.AutoMigrate(&models.User{}, &models.PasswordReset{}); err != nil {
		panic("couldn't migrate models to database! ")
	}
}
