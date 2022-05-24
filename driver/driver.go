package driver

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type MySQLConfig struct {
	Host     string
	User     string
	Password string
	Port     string
	Db       string
	Net      string
}


func ConnectToMySQL(conf MySQLConfig) (*sql.DB, error) {
	connectionString := fmt.Sprintf("%v:%v@%v(%v:%v)/%v", conf.User, conf.Password, conf.Net, conf.Host, conf.Port, conf.Db)

	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		return nil, err
	}

	return db, nil
}
