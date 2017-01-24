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
	SCHEMA_NAME = "price"
	PRICE_UPDATE_REQUEST_TABLE_NAME = "priceUpdateRequest"
	PRICE_TABLE_NAME = "price"
)

var db *sql.DB;

func Init() {
	var err error;
	db, err = sql.Open(DB_DRIVER, DB_CONNECTION)
	if (err != nil) {
		log.Fatal(err.Error())
	}
	db.Ping()
}

func GetDb() *sql.DB {
	return db;
}

func CloseDb() {
	db.Close()
}

func SavePriceInUpdateTable(priceObj *model.Product) {
	tx, err := db.Begin();
	if err != nil {
		log.Fatalf("Unable to create transection, error while saving update price request to db\n Error %v", err);
	}
	fmt.Println("Inserting record into priceUpdateRequest table...")
	insertQuery := fmt.Sprintf("insert into %s.%s (product_id, product_name,cost,status) values(%s,%s,%s,%s)", SCHEMA_NAME,
		PRICE_UPDATE_REQUEST_TABLE_NAME,
		priceObj.Product_id,
		priceObj.Product_name,
		priceObj.Cost,
		priceObj.Status)

	_, err = db.Exec(insertQuery)

	if (err != nil) {
		tx.Rollback();
		log.Fatalf("Unable to save update entry in db \nError: %v", err.Error());
	} else {
		tx.Commit();
	}
}

func ValidateEntryPresence(id int) {
	selectRecordCountByIdQuery := fmt.Sprintf("select count(*) from %s.%s where product_id = %s", SCHEMA_NAME, PRICE_TABLE_NAME, id);
	record, err := db.Query(selectRecordCountByIdQuery)
	if (err != nil) {
		log.Fatalf("Query failed while validating product id \n err: %v", err)
	}
	next := record.Next()

	if !next {
		log.Fatal("No record present to update")
	}
}

func GetUpdateRequestEntriesWithStatus(status string, tx sql.Tx) *sql.Rows {
	selectQuery := fmt.Sprintf("select product_id, product_name, cost from %s.%s where status='%s'", SCHEMA_NAME, PRICE_UPDATE_REQUEST_TABLE_NAME, status);
	fmt.Println(selectQuery);
	entries, err := tx.Query(selectQuery)
	if (err != nil) {
		log.Fatalf("Query failed while selecting update request \n err: %v", err)
	}
	return entries;
}

func ChangeStatusFrom(oldStatus string, newStatus string, tx sql.Tx) error {
	updateQuery := fmt.Sprintf("update %s.%s set status = '%s' where status = '%s'", SCHEMA_NAME, PRICE_UPDATE_REQUEST_TABLE_NAME, newStatus, oldStatus)

	_, err := tx.Exec(updateQuery)
	return err;
}

