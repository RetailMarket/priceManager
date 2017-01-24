package jobs

import (
	_"time"
	"log"
	"fmt"
	"Retail/priceManager/status"
	workflow "github.com/RetailMarket/workFlowClient"
	"golang.org/x/net/context"
	"Retail/priceManager/database"
)

func SendUpdatePriceForApprovalJob(client workflow.WorkFlowClient) {
	sendPendingRequestForApproval(client)
	//time.Sleep(time.Second * 10000)
}

func sendPendingRequestForApproval(client workflow.WorkFlowClient) {
	tx, err := database.GetDb().Begin()
	if err != nil {
		log.Fatalf("Begining transection failed while sending pending requests for approval  /n Errror: %v", err)
	}

	//for range time.Tick(time.Second * 2) {

	fmt.Println("Fetching pending update requests...")

	updatedValues := database.GetUpdateRequestEntriesWithStatus(status.PENDING, *tx)
	defer updatedValues.Close()

	request := &workflow.PriceUpdateRequest{}
	for updatedValues.Next() {
		fmt.Println("Got request :-)")
		var product_id int32
		var product_name string
		var cost float32

		updatedValues.Scan(&product_id, &product_name, &cost)

		priceObj := workflow.Product{
			ProductId: product_id,
			ProductName: product_name,
			Cost: cost,
			Status: status.PENDING}

		request.Products = append(request.Products, &priceObj)
	}

	response, err := client.SaveUpdatePriceForApproval(context.Background(), request)
	if err != nil {
		tx.Rollback()
		log.Printf("Failed to send pending update requests for approval \n err: %v\n", err)
	}

	err = database.ChangeStatusFrom(status.PENDING, status.PICKED, *tx)

	if (err != nil) {
		tx.Rollback()
		log.Printf("Query failed while changing status \n err: %v\n", err)
	}

	log.Printf("Response: %s", response.Message)
	tx.Commit();
	//}
}