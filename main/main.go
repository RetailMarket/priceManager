package main

import (
	"Retail/priceManager/database"
	"Retail/priceManager/seeds"
	"Retail/priceManager/server"
)

func main() {
	database.Init();
	defer database.CloseDb();
	seeds.UploadSeedForPriceTable();

	server.CreateServerConnection();
}