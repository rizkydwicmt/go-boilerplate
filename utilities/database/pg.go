package database

import (
	"fmt"
	authEntity "go-boilerplate/repository/auth/entity"
	helper "go-boilerplate/utilities"
	"strings"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnDB() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		viper.GetString("pg.host"),
		viper.GetInt("pg.port"),
		viper.GetString("pg.user"),
		viper.GetString("pg.password"),
		viper.GetString("pg.database"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if strings.ToLower(viper.GetString("ENVIRONMENT")) == "development" {
		db.AutoMigrate(&authEntity.User{})
	}

	sqlDB, err := db.DB()

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(viper.GetInt("pg.max_idle"))

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(viper.GetInt("pg.max_active"))

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err != nil {
		helper.Log("debug")("connection db", err)
	}
	return db
}
