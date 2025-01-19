package database

import (
	"database/sql"
	"fmt"
	"github.com/hoitek/Kit/retry"
	logger "github.com/hoitek/Logger"
	"github.com/hoitek/Maja-Service/config"
	_ "github.com/lib/pq"
)

// PostgresDB is a global variable for the raw sql connection
var PostgresDB *sql.DB

// ConnectPostgresDB connects to raw sql
func ConnectPostgresDB() *sql.DB {
	var (
		HOST     = config.AppConfig.DatabaseHost
		USER     = config.AppConfig.DatabaseUser
		PASSWORD = config.AppConfig.DatabasePassword
		DB_NAME  = config.AppConfig.DatabaseName
		PORT     = config.AppConfig.DatabasePort
		SSL_MODE = config.AppConfig.DatabaseSslMode
		TIMEZONE = config.AppConfig.DatabaseTimeZone
	)

	// Try to get raw sql connection for N amount of time
	conn, err := retry.Get(func() (*sql.DB, error) {
		DSN := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
			HOST,
			USER,
			PASSWORD,
			DB_NAME,
			PORT,
			SSL_MODE,
			TIMEZONE,
		)
		return sql.Open("postgres", DSN)
	}, 2, 3)

	// Handle error
	if err != nil {
		logger.Error(err)
		panic(err)
	} else {
		// Log when connection succeed
		logger.Info("Database connected!")
	}

	// Set global connection
	PostgresDB = conn

	return PostgresDB
}
