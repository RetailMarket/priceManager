package database

import (
	"database/sql"
	"Retail/priceManager/model"
	_ "github.com/bmizerany/pq"
	"fmt"
	"log"
	"strconv"
	"strings"
)

const (
	DB_DRIVER = "postgres"
	DB_CONNECTION = "user=postgres dbname=postgres password=postgres sslmode=disable"
	SCHEMA = "price"
	PRICE_UPDATE_REQUEST_TABLE = "priceUpdateRequest"
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

func getTx() *sql.Tx {
	tx, err := db.Begin();
	if err != nil {
		log.Fatalf("Unable to create transection \n Error %v", err);
	}
	return tx;
}

func joinIds(ids []int32) string {
	idsText := []string{}

	for i := range ids {
		number := ids[i]
		text := strconv.Itoa(int(number))
		idsText = append(idsText, text)
	}

	return strings.Join(idsText, "+")
}

func SavePriceInUpdateTable(priceObj *model.Product) {
	tx, err := db.Begin();
	if err != nil {
		log.Fatalf("Unable to create transection, error while saving update price request to db\n Error %v", err);
	}
	fmt.Println("Inserting record into priceUpdateRequest table...")
	insertQuery := fmt.Sprintf("insert into %s.%s (product_id, product_name,cost,status) values(%s,%s,%s,%s)", SCHEMA,
		PRICE_UPDATE_REQUEST_TABLE,
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
	selectQuery := fmt.Sprintf("select product_id, product_name, cost from %s.%s where status='%s'", SCHEMA, PRICE_UPDATE_REQUEST_TABLE, status);
	fmt.Println(selectQuery);
	entries, err := tx.Query(selectQuery)
	if (err != nil) {
		log.Fatalf("Query failed while selecting update request \n err: %v", err)
	}
	return entries;
}

func ChangeStatusFrom(oldStatus string, newStatus string, tx sql.Tx) error {
	updateQuery := fmt.Sprintf("update %s.%s set status = '%s' where status = '%s'", SCHEMA, PRICE_UPDATE_REQUEST_TABLE, newStatus, oldStatus)

	_, err := tx.Exec(updateQuery)
	return err;
}

func SwitchToLatest(ids []int32) error {
	tx := getTx();

	formattedIds := joinIds(ids);
	setCurrentVersionToNotLatestQuery := fmt.Sprintf("update %s.%s set is_latest = 'false' where product_id in (%s) and is_latest = true ", SCHEMA, PRICE_TABLE, formattedIds)

	_, err := tx.Exec(setCurrentVersionToNotLatestQuery);
	if err != nil {
		tx.Rollback();
		return err;
		log.Fatalf("Error while setting the current version to not latest \nError: %v", err);
	}
	setToLatestVersionQuery := fmt.Sprintf("update %s.%s set is_latest = 'true' from (select product_id pId, max(version) maxV from %s.%s group by product_id) as latestVersions where product_id = latestVersion.pid and version = latestVersions.maxV and product_id in %s", SCHEMA, PRICE_TABLE, SCHEMA, PRICE_TABLE, formattedIds);

	_, err = tx.Exec(setToLatestVersionQuery);

	if err != nil {
		tx.Rollback();
		return err;
		log.Fatalf("Error while setting the updated price to latest \nError: %v", err);
	}

	return nil;
}

func GetPriceUpdateRequests() (*sql.Rows, error) {
	selectQuery := fmt.Sprintf("select product_id, version from %s.%s inner join (select product_id pId, max(version) maxV from %s.%s group by product_id) latestVersions on product_id = latestVersions.pId and version = latestVersions.maxV where is_latest = false", SCHEMA, PRICE_TABLE, SCHEMA, PRICE_TABLE);

	fmt.Printf("Executing Query %s\n", selectQuery);
	return db.Query(selectQuery)
}