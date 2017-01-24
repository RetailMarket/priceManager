package main

import (
	"Retail/priceManager/database"
	"Retail/priceManager/client"
	"Retail/priceManager/seeds"
	"Retail/priceManager/jobs"
	_"Retail/priceManager/routes"
)

func main() {
	database.Init();
	defer database.CloseDb();

	// seeding data into DB.
	seeds.UploadSeedForPriceTable();
	seeds.UploadSeedForPriceUpdateRequestTable();

	workflowClient, conn := client.CreateClientConnection();

	// closing client connection.
	defer conn.Close();

	// running job for sending update price record for approval.
	jobs.SendUpdatePriceForApprovalJob(workflowClient)
	//routes.HandleRequest()
}