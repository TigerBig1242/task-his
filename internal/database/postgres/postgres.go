package postgres

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bigthamm/task-his/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDatabase(config *config.Config) (*gorm.DB, error) {
	// Connect to database
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DB_Host, config.DB_Port, config.DB_User, config.DB_Password,
		config.DB_Name)
	fmt.Println("DSN Connection string: ", dsn)

	// New logger for detailed SQL logging
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log Level
			Colorful:      true,        // Enable Color
		},
	)

	// Open Connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		return nil, fmt.Errorf("Failed to connect database: %w", err)
	}

	return db, nil
}
