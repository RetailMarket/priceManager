package database

import (
	"database/sql"
	_ "github.com/bmizerany/pq"
	priceClient "github.com/RetailMarket/priceManagerClient"
	"log"
	"Retail/priceManager/database/query"
	"Retail/priceManager/status"
	"strings"
	"strconv"
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

func AllLatestRecords() (*sql.Rows, error) {
	return db.Query(query.GetAllLatestRecordsQuery())
}

func SaveEntryForUpdate(tx *sql.Tx, req *priceClient.Record) error {
	var version string;
	var name string;
	err := tx.QueryRow(query.GetNewEntryDataQuery(int(req.ProductId))).Scan(&name, &version)
	if (err != nil) {
		log.Println("Unable to fetch record values for creating new update request")
		return err;
	}
	entry := &priceClient.Entry{
		ProductId: req.ProductId,
		ProductName: name,
		Version : getNextVersion(version),
		Cost: req.Cost,
		Status: status.PENDING,
		IsLatest: false,
	}
	_, err = tx.Exec(query.SaveNewRecordQuery(entry))
	return err;
}

func getNextVersion(currentVersion string) string {
	split := strings.Split(currentVersion, "")
	prefix := split[0]
	version, err := strconv.Atoi(strings.Join(split[1:], ""))
	if (err != nil) {
		log.Println(err)
	}
	return prefix + strconv.Itoa(version + VERSION_INC)
}