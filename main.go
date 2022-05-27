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

func dsn(dbName string) string {
	return f.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}

func dbConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn(""))
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return nil, err
	}
	// defer db.Close()

	// err = db.Ping()
	// if err != nil {
	// 	f.Println("error verifying connection with db.Ping")
	// 	panic(err.Error())
	// }
	// log.Printf("Verified connection from %s database \n", dbname)
	// return db, nil

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
		log.Printf("Error %s when opening database", err)
		return nil, err
	}
	//defer db.Close()

	// db.SetMaxOpenConns(20)
	// db.SetMaxIdleConns(20)
	// db.SetConnMaxLifetime(time.Minute * 5)

	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Errors %s pinging database", err)
		return nil, err
	}
	log.Printf("Verified connection from %s database with Ping\n", dbname)
	return db, nil
}

func insertToTable(db *sql.DB) error {
	insert, err := db.Query("INSERT INTO `RFID`.`Covert_RFID_JANCODE` (`drgm_rfid_cd`, `drgm_jan`, `drgm_jan2`) VALUES ('Test', 'Test', 'Test');")
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
	f.Println("Successful Insert to Database!")
	return nil
}

func convertFromRFID(db *sql.DB, rfid_code string) (string, string, error) {
	log.Printf("Getting JAN code")
	query := `select drgm_jan, drgm_jan2 from Covert_RFID_JANCODE where drgm_rfid_cd = ?;`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return "", "", err
	}
	defer stmt.Close()
	var JAN_code, JAN_code_2 string
	row := stmt.QueryRowContext(ctx, rfid_code)
	if err := row.Scan(&JAN_code, &JAN_code_2); err != nil {
		return "", "", err
	}
	return JAN_code, JAN_code_2, nil

}

// type book struct {
// 	rfid_code  string
// 	jan_code   string
// 	jan_code_2 string
// }

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

	rfid_code := "0xRFID"
	jan_code, jan_code_2, err := convertFromRFID(db, rfid_code)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("Product %s not found in DB", rfid_code)
	case err != nil:
		log.Printf("Encountered err %s when fetching price from DB", err)
	default:
		log.Printf("JanCode 1 of %s is %s, JanCode 2 is %s", rfid_code, jan_code, jan_code_2)
	}

	// for _, y := range x {
	// 	log.Printf("Name: %s Price: %d", x.name, x.price)
	// }

	// jancodes, err := selectJanCodeByRFID(db, rfid_code)
	// if err != nil {
	// 	log.Printf("Error %s when selecting product by price", err)
	// 	return
	// }
	// for _, book := range jancodes {
	// 	log.Printf("Name: %s Price: %s", book.rfid_code, book.jan_code)
	// }

	// err = db.QueryRow("select drgm_jan from Covert_RFID_JANCODE where drgm_rfid_cd = ?", 1).Scan(&rfid_code)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// f.Println(rfid_code)

}
