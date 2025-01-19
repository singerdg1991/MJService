package database

import (
	"testing"

	"github.com/hoitek/Maja-Service/config"
)

func TestConnectPostgresDB(t *testing.T) {
	config.LoadDefault()
	dbInstance := ConnectPostgresDB()
	if dbInstance == nil || PostgresDB == nil {
		t.Error("Connection by raw client to database has been faild.")
	}
}
