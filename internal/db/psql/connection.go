package psql

import (
	"fmt"
	"go-finance-tracker/internal/config"
	"go-finance-tracker/internal/models"
	"go-finance-tracker/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB = nil

func initDB(config config.PostgresDB) error {
	dbConnString := fmt.Sprintf(
		// example: postgres://postgres:postgres@localhost:5432/postgres?sslmode=false
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
		config.Sslmode,
	)

	var err error
	db, err = gorm.Open(postgres.Open(dbConnString), &gorm.Config{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.FinanceRecord{},
		&models.TransactionType{},
		&models.Category{},
	)

	if err != nil {
		logger.GetLogger().Error("❌ Auto migration failed")
		return err
	}
	logger.GetLogger().Info("👍 Migration complete - gorm service")
	return nil
}

func GetDbInstance(database config.PostgresDB) (*gorm.DB, error) {
	if db == nil {
		if err := initDB(database); err != nil {
			return nil, err
		}
	}

	return db, nil
}
