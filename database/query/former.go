package query

import (
	"fmt"
	priceClient "github.com/RetailMarket/priceManagerClient"
	"strconv"
	"strings"
	"Retail/priceManager/status"
)

const (
	SCHEMA = "price"
	PRICE_TABLE = "price"
)

func joinProductIds(records []*priceClient.Entry) string {
	idsText := []string{}

	for i := range records {
		number := records[i].ProductId
		text := strconv.Itoa(int(number))
		idsText = append(idsText, text)
	}

	return strings.Join(idsText, ",")
}

func ChangeStatusQuery(status string, product_id int, version string) string {
	return fmt.Sprintf("update %s.%s set status = '%s' where product_id = %d and version = '%s'", SCHEMA, PRICE_TABLE, status, int(product_id), version)
}

func SetNotLatestQuery(records []*priceClient.Entry) string {
	formattedIds := joinProductIds(records);
	return fmt.Sprintf("update %s.%s set is_latest = false where product_id in (%s) and is_latest = true", SCHEMA, PRICE_TABLE, formattedIds)
}

func SetToLatestQuery(product_id int, version string) string {
	return fmt.Sprintf("update %s.%s set is_latest = true where product_id = %d and version = '%s'", SCHEMA, PRICE_TABLE, product_id, version);
}

func GetPendingRecordsQuery() string {
	return fmt.Sprintf("select product_id, version from %s.%s inner join (select product_id pId, max(version) maxV from %s.%s group by product_id) latestVersions on product_id = latestVersions.pId and version = latestVersions.maxV where is_latest = false and status = '%s'", SCHEMA, PRICE_TABLE, SCHEMA, PRICE_TABLE, status.PENDING);
}

func GetAllLatestRecordsQuery() string {
	return fmt.Sprintf("select * from %s.%s where is_latest = true", SCHEMA, PRICE_TABLE);
}

func SaveNewRecordQuery(entry *priceClient.Entry) string {
	return fmt.Sprintf("insert into %s.%s (product_id,product_name,cost,version,status,is_latest) values (%d,'%s',%d,'%s','%s',%t)", SCHEMA, PRICE_TABLE, entry.ProductId, entry.ProductName, entry.Cost, entry.Version, entry.Status, entry.IsLatest);
}

func GetNewEntryDataQuery(id int) string {
	return fmt.Sprintf("select product_name, max(version) from %s.%s where product_id = %d group by product_name", SCHEMA, PRICE_TABLE, id);
}