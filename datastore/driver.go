package datastore

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type mySQLConfig struct {
	Host     string
	User     string
	Password string
	Port     string
	Db       string
	Net      string
}

func ConnectToMySQL() (*sql.DB, error) {
		
	envErr := godotenv.Load()

	if envErr != nil {
		log.Fatal("Error loading .env file")
	}

	var conf = mySQLConfig{
		Host:     os.Getenv("HOST"),
		User:     os.Getenv("DBUSER"),
		Password: os.Getenv("DBPASS"),
		Port:     os.Getenv("PORT"),
		Db:       os.Getenv("DB"),
		Net:      "tcp",
	}

	connectionString := fmt.Sprintf("%v:%v@%v(%v:%v)/%v", conf.User, conf.Password, conf.Net, conf.Host, conf.Port, conf.Db)

	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		return nil, err
	}

	return db, nil
}
