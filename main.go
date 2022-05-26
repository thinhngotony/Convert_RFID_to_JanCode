package main

import (
	"context"
	"database/sql"
	f "fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	username = "root"
	password = "Rfid@2022"
	hostname = "192.168.1.244:3306"
	dbname   = "RFID"
)

func insertToTable(db *sql.DB) error {

	insert, err := db.Query("INSERT INTO `RFID`.`Covert_RFID_JANCODE` (`drgm_rfid_cd`, `drgm_jan`, `drgm_jan2`) VALUES ('Test', 'Test', 'Test');")
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
	f.Println("Successful Insert to Database!")

	return nil
}

func dsn(dbName string) string {
	return f.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}

func dbConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn(""))
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return nil, err
	}
	//defer db.Close()

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+dbname)
	if err != nil {
		log.Printf("Error %s when creating DB\n", err)
		return nil, err
	}
	no, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when fetching rows", err)
		return nil, err
	}
	log.Printf("rows affected %d\n", no)

	db.Close()
	db, err = sql.Open("mysql", dsn(dbname))
	if err != nil {
		log.Printf("Error %s when opening DB", err)
		return nil, err
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Errors %s pinging DB", err)
		return nil, err
	}
	log.Printf("Connected to database %s successfully\n", dbname)
	return db, nil
}

func main() {

	db, err := dbConnection()
	if err != nil {
		log.Printf("Error %s when getting db connection", err)
		return
	}
	defer db.Close()
	log.Printf("Successfully connected to database")
	// err = insertToTable(db)
	// if err != nil {
	// 	log.Printf("Insert failed failed with error %s", err)
	// 	return
	// }

}
