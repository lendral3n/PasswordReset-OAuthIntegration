package database

import (
	"fmt"
	"emailnotifl3n/app/config"
	ud "emailnotifl3n/features/user/data"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDBPostgres(cfg *config.AppConfig) *gorm.DB {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable TimeZone=Asia/Shanghai",
		cfg.DB_HOSTNAME, cfg.DB_PORT, cfg.DB_USERNAME, cfg.DB_NAME, cfg.DB_PASSWORD)

	DB, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	DB.AutoMigrate(
		&ud.User{},
	)

	return DB
}