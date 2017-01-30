package seeds

import (
	"fmt"
	"Retail/priceManager/database"
	"Retail/priceManager/status"
	"log"
)

type Product struct {
	product_id   int
	product_name string
	cost         int
	version      string
	is_latest    bool
	status       string
}

type Config struct {
	Products []Product
}

func UploadSeedForPriceTable() {
	db := database.GetDb()
	fmt.Println("Seeding data to price table...")

	var config Config

	p1 := Product{product_id: 1001, product_name: "Pen", cost: 10, version:"v1", status: status.COMPLETED, is_latest:false}
	p2 := Product{product_id: 1001, product_name: "Pen", cost: 12, version:"v2", status: status.COMPLETED, is_latest:true}
	p3 := Product{product_id: 1001, product_name: "Pen", cost: 14, version:"v3", status: status.PENDING, is_latest:false}

	p4 := Product{product_id: 1002, product_name: "Pencil", cost: 4, version:"v1", status: status.COMPLETED, is_latest:true}
	p5 := Product{product_id: 1002, product_name: "Pencil", cost: 6, version:"v2", status: status.PENDING, is_latest:false}

	p6 := Product{product_id: 1003, product_name: "Rubber", cost: 2, version:"v1", status: status.COMPLETED, is_latest:false}
	p7 := Product{product_id: 1003, product_name: "Rubber", cost: 3, version:"v2", status: status.COMPLETED, is_latest:false}
	p8 := Product{product_id: 1003, product_name: "Rubber", cost: 5, version:"v3", status: status.COMPLETED, is_latest:true}
	p9 := Product{product_id: 1003, product_name: "Rubber", cost: 8, version:"v4", status: status.PENDING, is_latest:false}

	p10 := Product{product_id: 1004, product_name: "Sticky", cost: 10, version:"v1", status: status.COMPLETED, is_latest:false}
	p11 := Product{product_id: 1004, product_name: "Sticky", cost: 13, version:"v2", status: status.COMPLETED, is_latest:true}

	p12 := Product{product_id: 1005, product_name: "chart", cost: 5, version:"v1", status: status.COMPLETED, is_latest:true}

	p13 := Product{product_id: 1006, product_name: "global map", cost: 2, version:"v1", status: status.COMPLETED, is_latest:false}
	p14 := Product{product_id: 1006, product_name: "global map", cost: 3, version:"v2", status: status.COMPLETED, is_latest:true}
	p15 := Product{product_id: 1006, product_name: "global map", cost: 4, version:"v3", status: status.PENDING, is_latest:false}

	config.Products = append(config.Products, p1, p2, p3, p4, p5, p6, p7, p8, p9, p10, p11, p12, p13, p14, p15)

	_, err := db.Exec("truncate table price.price")
	if (err != nil) {
		panic(err.Error())
	}

	for i := 0; i < len(config.Products); i++ {
		product := config.Products[i];
		insertQuery := fmt.Sprintf("insert into price.price (product_id,product_name,cost,version,status,is_latest) values (%d,'%s',%d,'%s','%s',%t)", product.product_id, product.product_name, product.cost, product.version, product.status, product.is_latest)
		log.Println(insertQuery)
		_, err := db.Exec(insertQuery)
		if (err != nil) {
			log.Printf("Error %v", err)
			return;
		}
	}
	fmt.Println("Seed successful for price table :)")
}

