package seeds

import (
	"fmt"
	"Retail/priceManager/status"
	"Retail/priceManager/database"
)

type Product struct {
	product_id   int
	product_name string
	cost         int
	status       string
}

type Config struct {
	Products []Product
}

func UploadSeedForPriceTable() {
	db := database.GetDb()
	fmt.Println("Seeding data to price table...")

	var config Config

	p1 := Product{product_id: 1001, product_name: "Pen", cost: 10}
	p2 := Product{product_id: 1002, product_name: "Pencil", cost: 4}
	p3 := Product{product_id: 1003, product_name: "Rubber", cost: 2}
	p4 := Product{product_id: 1004, product_name: "Sticky", cost: 10}
	p5 := Product{product_id: 1005, product_name: "chart", cost: 5}
	p6 := Product{product_id: 1006, product_name: "global map", cost: 2}

	config.Products = append(config.Products, p1, p2, p3, p4, p5, p6)

	_, err := db.Exec("truncate table price.price")
	if (err != nil) {
		panic(err.Error())
	}

	for i := 0; i < len(config.Products); i++ {
		product := config.Products[i];
		_, err := db.Exec("insert into price.price (product_id,product_name,cost) values ($1,$2,$3)", product.product_id, product.product_name, product.cost)
		if (err != nil) {
			panic(err)
		}
	}
	fmt.Println("Seed successful for price table :)")
}

func UploadSeedForPriceUpdateRequestTable() {
	db := database.GetDb()
	fmt.Println("Seeding data to price update table...")

	var config Config

	p1 := Product{product_id: 1001, product_name: "Pen", cost: 12, status: status.PENDING}
	p2 := Product{product_id: 1002, product_name: "Pencil", cost: 5, status: status.PENDING}
	p3 := Product{product_id: 1003, product_name: "Rubber", cost: 3, status: status.PENDING}
	p4 := Product{product_id: 1004, product_name: "Sticky", cost: 20, status: status.COMPLETED}
	p5 := Product{product_id: 1005, product_name: "chart", cost: 10, status: status.COMPLETED}

	config.Products = append(config.Products, p1, p2, p3, p4, p5)

	_, err := db.Exec("truncate table price.priceUpdateRequest")
	if (err != nil) {
		panic(err.Error())
	}

	for i := 0; i < len(config.Products); i++ {
		product := config.Products[i];
		_, err := db.Exec("insert into price.priceUpdateRequest (product_id,product_name,cost,status) values ($1,$2,$3,$4)", product.product_id, product.product_name, product.cost, product.status)
		if (err != nil) {
			panic(err)
		}

	}
	fmt.Println("Seed successful for price update table :)")
}

