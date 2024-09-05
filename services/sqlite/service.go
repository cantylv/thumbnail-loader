package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"go.uber.org/zap"
)

func Init(logger *zap.Logger) *sql.DB {
	db, err := sql.Open("sqlite3", "./services/sqlite/data/database.db")
	if err != nil {
		logger.Fatal(fmt.Sprintf("fatal error while connecting to sqlite: %v", err))
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)
	db.SetConnMaxIdleTime(time.Minute * 1)

	logger.Info("succesful connection to Sqlite")
	return db
}
