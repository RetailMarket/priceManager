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
	"Retail/priceManager/status"
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

func (s *server) PriceUpdateRecords(ctx context.Context, _ *priceClient.FetchRecordsRequest) (*priceClient.FetchRecordsResponse, error) {
	records, err := database.PriceUpdateRequests();
	response := &priceClient.FetchRecordsResponse{}
	if (err != nil) {
		log.Printf("Query failed while selecting update request \n err: %v", err)
	} else {
		for records.Next() {
			var product_id int32;
			var version string;
			records.Scan(&product_id, &version)
			record := priceClient.Entry{
				ProductId: product_id,
				Version: version}
			response.Entries = append(response.Entries, &record)
		}
	}
	return response, err
}

func (s *server) NotifySuccessfullyPicked(ctx context.Context, request *priceClient.NotifyRequest) (*priceClient.NotifyResponse, error) {
	records := request.GetEntries();
	response := &priceClient.NotifyResponse{}

	tx, err := database.GetDb().Begin();
	if (err != nil) {
		response.Message = fmt.Sprintln("Error while creating database transection");
		return response, err;
	}

	err = database.ChangeStatusTo(tx, status.PICKED, records);

	if (err != nil) {
		tx.Rollback();
		response.Message = fmt.Sprintf("Failed to change status of %v to picked", records);
	} else {
		tx.Commit();
		response.Message = fmt.Sprintf("Successfully changed status of %v to picked", records);
	}

	return response, err
}

func (s *server) NotifySuccessfullyProcessed(ctx context.Context, request *priceClient.NotifyRequest) (*priceClient.NotifyResponse, error) {
	records := request.GetEntries();
	tx, err := database.GetDb().Begin();
	message := fmt.Sprintf("Successfully changed status to completed and set %v to latest", records);
	response := &priceClient.NotifyResponse{Message:message}

	if (err != nil) {
		response.Message = fmt.Sprintln("Error while creating database transection");
		return response, err;
	}

	err = database.ChangeStatusTo(tx, status.COMPLETED, records);

	if (err != nil) {
		tx.Rollback();
		response.Message = fmt.Sprintf("Failed to change status of %v to picked", records);
		return response, err
	} else {
		err = database.SwitchToLatest(tx, records)
		if (err != nil) {
			tx.Rollback()
			response.Message = fmt.Sprintf("unable to set %v to latest", records)
			return response, err
		}
	}

	tx.Commit();
	return response, err
}

func (s *server) AllRecords(ctx context.Context, _ *priceClient.FetchRecordsRequest) (*priceClient.FetchRecordsResponse, error) {
	records, err := database.AllRecords();
	response := &priceClient.FetchRecordsResponse{}
	if (err != nil) {
		log.Printf("Query failed while selecting all records \n err: %v", err)
	} else {
		for records.Next() {
			var product_id int32;
			var version string;
			records.Scan(&product_id, &version)
			record := priceClient.Entry{
				ProductId: product_id,
				Version: version}
			response.Entries = append(response.Entries, &record)
		}
	}
	return response, err
}