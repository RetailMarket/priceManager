package jobs

import (
	"time"
	"database/sql"
	"log"
	"Retail/priceManager/model"
	"fmt"
	"Retail/priceManager/status"
)

func SendUpdatePriceForApprovalJob(db *sql.DB) {
	go updateNewEntries(db)
	time.Sleep(time.Second * 10000)
}
func updateNewEntries(db *sql.DB) {
	for range time.Tick(time.Second * 2) {
		fmt.Println("fetching values...")
		selectQuery := "select product_id, product_name, cost from price.updatePriceRequest where status = $1"
		//updateQuery := "update price.update_price set status = 'PICKED' where status = 'PENDING'"
		updatedValues, err := db.Query(selectQuery, status.PENDING)
		if (err != nil) {
			log.Fatal(err.Error())
		} else {
			fmt.Println("No Error...")
			listOfUpdatedRow := []*model.Price{}

			for updatedValues.Next() {
				fmt.Println("Got values :-)")
				var product_id int
				var product_name string
				var cost float64

				updatedValues.Scan(&product_id, &product_name, &cost)

				priceObj := model.Price{
					Product_id: product_id,
					Product_name: product_name,
					Cost: cost,
					Status: status.PENDING}

				listOfUpdatedRow = append(listOfUpdatedRow, &priceObj)
			}

			for i := 0; i < len(listOfUpdatedRow); i++ {
				fmt.Println(*listOfUpdatedRow[i])
			}
		}
	}
}