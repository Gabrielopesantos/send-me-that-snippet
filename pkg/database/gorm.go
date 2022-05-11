package database

import (
	"fmt"
	"github.com/gabrielopesantos/smts/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGormDB(cfg *config.Config) (*gorm.DB, error) {
	driver := postgres.Open
	//pgDsn := "postgres://gabriel:gabriel@localhost:5432/main"
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.DBConfig.User, cfg.DBConfig.Password, cfg.DBConfig.Host, cfg.DBConfig.Port,
		cfg.DBConfig.Database)

	db, err := gorm.Open(driver(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
