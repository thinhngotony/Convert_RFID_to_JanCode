package main

import (
	"database/sql"
	f "fmt"
	"log"

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

	err = db.Ping()
	if err != nil {
		f.Println("error verifying connection with db.Ping")
		panic(err.Error())
	}
	log.Printf("Verified connection from %s database \n", dbname)
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

	// Structure to QUERY with MySQL database//
	//========================================//
	// err = insertToTable(db)
	// if err != nil {
	// 	log.Printf("Insert failed failed with error %s", err)
	// 	return
	// }
	//========================================//

}
