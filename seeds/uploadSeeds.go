package seeds

import (
	"fmt"
	"database/sql"
)

type Product struct {
	product_id   int
	product_name string
	cost         int
	//status       string
}

type Config struct {
	Products []Product
}

func UploadSeedData(db *sql.DB) {
	fmt.Println("Seeding data...")

	var config Config

	p1 := Product{product_id: 1001, product_name: "Pen", cost: 12}
	p2 := Product{product_id: 1002, product_name: "Pencil", cost: 5}
	p3 := Product{product_id: 1003, product_name: "Rubber", cost: 3}
	p4 := Product{product_id: 1004, product_name: "Sticky", cost: 20}
	p5 := Product{product_id: 1005, product_name: "chart", cost: 10}
	p6 := Product{product_id: 1006, product_name: "global map", cost: 2}

	config.Products = append(config.Products, p1, p2, p3, p4, p5, p6)

	_, err := db.Exec("truncate table price.price")
	if (err != nil) {
		panic(err)
	}

	for i := 0; i < len(config.Products); i++ {
		product := config.Products[i];
		_, err := db.Exec("INSERT INTO price.price (product_id,product_name,cost) values ($1,$2,$3)", product.product_id, product.product_name, product.cost)
		if (err != nil) {
			panic(err)
		}
	}
	fmt.Println("Seed successful :)")
}