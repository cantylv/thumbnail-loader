package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/cantylv/thumbnail-loader/services/connectors"
	_ "github.com/mattn/go-sqlite3"

	"go.uber.org/zap"
)

var _ connectors.ClientDB = (*sql.DB)(nil)

type ClientInstance struct{}

var _ connectors.EngineDB = (*ClientInstance)(nil)

func NewClientInstanse() *ClientInstance {
	return &ClientInstance{}
}

func (t ClientInstance) InitClientDB(logger *zap.Logger) connectors.ClientDB {
	db, err := sql.Open("sqlite3", "./services/sqlite/data/database.db")
	if err != nil {
		logger.Panic(fmt.Sprintf("fatal error while connecting to sqlite: %v", err))
	}
	db.SetConnMaxLifetime(time.Minute * 3)

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)
	db.SetConnMaxIdleTime(time.Minute * 1)

	logger.Info("succesful connection to Sqlite")
	return db
}
