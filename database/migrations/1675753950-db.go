package migrations

import (
	"database/sql"
	"fmt"
	"github.com/hoitek/Kit/retry"
	"github.com/hoitek/Maja-Service/config"
)

func MigrateDB() error {
	// Connect to the PostgresSQL server
	var (
		user     = config.AppConfig.DatabaseUser
		password = config.AppConfig.DatabasePassword
		host     = config.AppConfig.DatabaseHost
		port     = config.AppConfig.DatabasePort
		sslMode  = config.AppConfig.DatabaseSslMode
		dbName   = config.AppConfig.DatabaseName
	)

	// Connect to the database
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/?sslmode=%s", user, password, host, port, sslMode)
	db, err := retry.Get(func() (*sql.DB, error) {
		return sql.Open("postgres", connStr)
	}, 2, 3)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Check if the database exists
	var existsDB bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", dbName).Scan(&existsDB)
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %v", err)
	}

	// Create the database if it doesn't exist
	if !existsDB {
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s OWNER %s", dbName, user))
		if err != nil {
			return fmt.Errorf("failed to create database: %v", err)
		}
	}

	// Check if the user exists
	var existsUser bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_roles WHERE rolname = $1)", user).Scan(&existsUser)
	if err != nil {
		return fmt.Errorf("failed to check if user exists: %v", err)
	}

	// Create the user if it doesn't exist
	if !existsUser {
		_, err = db.Exec(fmt.Sprintf("CREATE USER %s WITH PASSWORD '%s'", user, password))
		if err != nil {
			return fmt.Errorf("failed to create user: %v", err)
		}
	}

	// Grant privileges to the user on the database
	db.Exec(fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE %s TO %s", dbName, user))

	return nil
}
