package database

import (
	"database/sql"
	"Retail/priceManager/model"
	_ "github.com/bmizerany/pq"
	priceClient "github.com/RetailMarket/priceManagerClient"
	"fmt"
	"log"
	"strconv"
	"strings"
	"Retail/priceManager/status"
)

const (
	DB_DRIVER = "postgres"
	DB_CONNECTION = "user=postgres dbname=postgres password=postgres sslmode=disable"
	SCHEMA = "price"
	PRICE_TABLE = "price"
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

func joinProductIds(records []*priceClient.Entry) string {
	idsText := []string{}

	for i := range records {
		number := records[i].ProductId
		text := strconv.Itoa(int(number))
		idsText = append(idsText, text)
	}

	return strings.Join(idsText, ",")
}

func SavePriceInUpdateTable(priceObj *model.Product) {
	tx, err := db.Begin();
	if err != nil {
		log.Fatalf("Unable to create transection, error while saving update price request to db\n Error %v", err);
	}
	fmt.Println("Inserting record into priceUpdateRequest table...")
	insertQuery := fmt.Sprintf("insert into %s.%s (product_id, product_name,cost,status) values(%s,%s,%s,%s)", SCHEMA,
		PRICE_TABLE,
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
	selectRecordCountByIdQuery := fmt.Sprintf("select count(*) from %s.%s where product_id = %s", SCHEMA, PRICE_TABLE, id);
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
	selectQuery := fmt.Sprintf("select product_id, product_name, cost from %s.%s where status='%s'", SCHEMA, PRICE_TABLE, status);
	fmt.Println(selectQuery);
	entries, err := tx.Query(selectQuery)
	if (err != nil) {
		log.Fatalf("Query failed while selecting update request \n err: %v", err)
	}
	return entries;
}

func ChangeStatusTo(tx *sql.Tx, status string, records []*priceClient.Entry) error {
	for i := 0; i < len(records); i++ {
		updateQuery := fmt.Sprintf("update %s.%s set status = '%s' where product_id = %d and version = '%s'", SCHEMA, PRICE_TABLE, status, int(records[i].ProductId), records[i].Version)
		_, err := tx.Exec(updateQuery)
		if (err != nil) {
			return err;
		}
	}
	return nil;
}

func SwitchToLatest(tx *sql.Tx, records []*priceClient.Entry) error {
	formattedIds := joinProductIds(records);

	setCurrentVersionToNotLatestQuery := fmt.Sprintf("update %s.%s set is_latest = false where product_id in (%s) and is_latest = true ", SCHEMA, PRICE_TABLE, formattedIds)

	_, err := tx.Exec(setCurrentVersionToNotLatestQuery);

	for i := range records {
		setToLatestVersionQuery := fmt.Sprintf("update %s.%s set is_latest = true where product_id = %d and version = '%s'", SCHEMA, PRICE_TABLE, int(records[i].ProductId), records[i].Version);
		_, err = tx.Exec(setToLatestVersionQuery);
		if (err != nil) {
			return err;
		}
	}
	//setToLatestVersionQuery := fmt.Sprintf("update %s.%s set is_latest = 'true' from (select product_id pId, max(version) maxV from %s.%s group by product_id) as latestVersions where product_id = latestVersion.pId and version = latestVersions.maxV and product_id in (%s)", SCHEMA, PRICE_TABLE, SCHEMA, PRICE_TABLE, formattedIds);
	//fmt.Print(setToLatestVersionQuery)

	return err;
}

func GetPriceUpdateRequests() (*sql.Rows, error) {
	selectQuery := fmt.Sprintf("select product_id, version from %s.%s inner join (select product_id pId, max(version) maxV from %s.%s group by product_id) latestVersions on product_id = latestVersions.pId and version = latestVersions.maxV where is_latest = false and status = '%s'", SCHEMA, PRICE_TABLE, SCHEMA, PRICE_TABLE, status.PENDING);

	fmt.Printf("Executing Query %s\n", selectQuery);
	return db.Query(selectQuery)
}