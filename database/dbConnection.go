package database

import (
	"database/sql"
	"Retail/priceManager/model"
	_ "github.com/bmizerany/pq"
	"fmt"
	"log"
)

const (
	DB_DRIVER = "postgres"
	DB_CONNECTION = "user=postgres dbname=postgres password=postgres sslmode=disable"
)

func OpenDatabase() *sql.DB {
	db, err := sql.Open(DB_DRIVER, DB_CONNECTION)
	if (err != nil) {
		log.Fatal(err.Error())
	}
	db.Ping()
	return db;
}

func SavePriceInUpdateTable(db *sql.DB, priceObj *model.Price) {
	fmt.Println("Inserting record into updatePriceRequest table...")
	insertQuery := "insert into price.updatePriceRequest (product_id, product_name,cost,status) values($1,$2,$3,$4)"

	_, err := db.Exec(
		insertQuery,
		priceObj.Product_id,
		priceObj.Product_name,
		priceObj.Cost,
		priceObj.Status)

	if (err != nil) {
		log.Fatalf("Unable to save update entry in db \nError: %v", err.Error());
	}
}

func ValidateEntryPresence(db *sql.DB, id int) {
	selectRecordCountByIdQuery := "select count(*) from price.price where product_id = $1"
	record, err := db.Query(selectRecordCountByIdQuery, id)
	if(err != nil){
		log.Fatal(err)
	}

	next := record.Next()

	if !next {
		log.Fatal("No record present to update")
	}
}

func CleanUpUpdatePriceRequestTable(db *sql.DB) {
	truncateQuery := "truncate table price.updatePriceRequest"
	_, err := db.Exec(truncateQuery)
	if (err != nil) {
		log.Fatal("Failed to clean uo update price request table")
	}
}