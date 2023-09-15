package database

import (
	"database/sql"
	"fmt"

	"github.com/renan5g/go-clean-arch/config"
)

func Open(configs *config.Conf) *sql.DB {
	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	return db
}
