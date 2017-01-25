package server

import (
	"net"
	"log"
	"fmt"
	"google.golang.org/grpc"
	priceClient "github.com/RetailMarket/priceManagerClient"
	"google.golang.org/grpc/reflection"
	"Retail/priceManager/database"
	"golang.org/x/net/context"
)

const (
	port = ":3000"
)

type server struct{}

func CreateServerConnection() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	_server := grpc.NewServer()

	priceClient.RegisterPriceManagerServer(_server, &server{});

	reflection.Register(_server)
	fmt.Printf("Price Server Listening to Port: %s\n", port);
	if err := _server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) SwitchToLatestVersion(ctx context.Context, idContainer *priceClient.ChangeLatestRequest) (*priceClient.ChangeLatestResponse, error) {
	ids := idContainer.GetProductId();
	err := database.SwitchToLatest(ids);
	message := "";
	if (err != nil) {
		message = fmt.Sprintf("Failed to switch set new prices to latest Ids: %v", ids);
	}
	message = fmt.Sprintf("Successfully set the new prices for Ids: %v", ids);
	return &priceClient.ChangeLatestResponse{Message: message}, err
}

func (s *server) GetPriceUpdateRecords(ctx context.Context, idContainer *priceClient.FetchRecordsRequest) (*priceClient.FetchRecordsResponse, error) {
	records, err := database.GetPriceUpdateRequests();
	response := &priceClient.FetchRecordsResponse{}
	if (err != nil) {
		log.Fatalf("Query failed while selecting update request \n err: %v", err)
	} else {
		for records.Next() {
			var product_id int32;
			var version string;
			records.Scan(&product_id, &version)
			record := priceClient.UpdateProductPriceEntry{
				ProductId: product_id,
				Version: version}
			response.Products = append(response.Products, &record)
		}
	}
	return response, err
}
