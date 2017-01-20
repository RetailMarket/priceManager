package main

import (
	"Retail/priceManager/routes"
	"Retail/priceManager/database"
	//"Retail/priceManager/jobs"
	"Retail/priceManager/seeds"
)

func main() {
	db := database.OpenDatabase()
	defer db.Close()

	// for seeding the price table
	seeds.UploadSeedData(db)
	database.CleanUpUpdatePriceRequestTable(db)

	//jobs.SendUpdatePriceForApprovalJob(db)
	routes.HandleRequest(db)
}