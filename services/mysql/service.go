package mysql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func Init(logger *zap.Logger) *sql.DB {
	address := fmt.Sprintf("%s:%d", viper.GetString("mysql.host"), viper.GetUint16("mysql.port"))
	connLine := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		address,
		viper.GetString("mysql.dbname"))
	db, err := sql.Open("mysql", connLine)
	if err != nil {
		logger.Fatal(fmt.Sprintf("fatal error while connecting to mysql: %v", err))
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)
	db.SetConnMaxIdleTime(time.Minute * 1)

	logger.Info("succesful connection to Mysql")
	return db
}
