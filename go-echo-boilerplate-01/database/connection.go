package database

import (
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

type DatabaseProvider interface{}

type DatabaseProviderConnection struct {
	Db *sqlx.DB
}

func ConnectMYSQL() DatabaseProvider {
	// Connecting to MYSQL database
	config := map[string]string{
		"user":     viper.GetString("database.user"),
		"password": viper.GetString("database.pwd"),
		"host":     viper.GetString("database.host"),
		"port":     viper.GetString("database.port"),
		"db_name":  viper.GetString("database.db_name"),
	}

	db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config["user"],
		config["password"],
		config["host"],
		config["port"],
		config["db_name"],
	))

	if err != nil {
		log.Fatalf("Could not connect to database :%v", err)
	} else {
		log.Println("Database Connected!")
	}

	// db.MapperFunc(strings.ToLower)
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return &DatabaseProviderConnection{
		Db: db,
	}
}
