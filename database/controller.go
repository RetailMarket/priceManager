package database

import (
	"database/sql"
	_ "github.com/bmizerany/pq"
	priceClient "github.com/RetailMarket/priceManagerClient"
	"log"
	"Retail/priceManager/database/query"
)

const (
	DB_DRIVER = "postgres"
	DB_CONNECTION = "user=postgres dbname=postgres password=postgres sslmode=disable"
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

func ChangeStatusTo(tx *sql.Tx, status string, records []*priceClient.Entry) error {
	var err error;
	for i := 0; i < len(records); i++ {
		product_id := int(records[i].ProductId)
		version := records[i].Version
		updateQuery := query.ChangeStatusQuery(status, product_id, version)
		_, err = tx.Exec(updateQuery)
		log.Printf("Changing status of product id %d version %s :Error %v\n", product_id, version, err)
	}
	return err;
}

func SwitchToLatest(tx *sql.Tx, records []*priceClient.Entry) error {
	_, err := tx.Exec(query.SetNotLatestQuery(records));
	log.Printf("Executing query to set flag is_latest to false for entries %v :Error %v\n", records, err);

	for i := range records {
		_, err = tx.Exec(query.SetToLatestQuery(int(records[i].ProductId), records[i].Version));
		log.Printf("Executing query to set flag is_latest to true for entry %v :Error %v\n", records[i], err);
	}

	return err;
}

func PriceUpdateRequests() (*sql.Rows, error) {
	return db.Query(query.GetPendingRecordsQuery())
}

func AllRecords() (*sql.Rows, error) {
	return db.Query(query.GetAllRecordsQuery())
}