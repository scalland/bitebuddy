package utils

import (
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
)

func (u *Utils) ConnectDB() (*sql.DB, error) {
	// Connect to the database (adjust DSN as needed)
	// DSN format: username:password@tcp(host:port)/dbname?parseTime=true
	return sql.Open(viper.GetString("db_driver"), fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", viper.GetString("db_username"), viper.GetString("db_password"), viper.GetString("db_host"), viper.GetInt("db_port"), viper.GetString("db_database")))
}

func (u *Utils) CloseDB(db *sql.DB) error {
	if db != nil {
		return db.Close()
	}
	return nil
}
