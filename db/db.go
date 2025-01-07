package db

import (
	"fmt"
	"learn/settings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBManager struct {
	Dsn string
}

var Database *gorm.DB

func init() {
	db, err := gorm.Open(postgres.Open(settings.Settings.Postgres.PostgresDsn()), &gorm.Config{})
	if err != nil {
		fmt.Println("Error db conn")
		panic(err)
	}
	Database = db
}
