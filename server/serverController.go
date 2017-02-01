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

func (s *server) PendingRecords(ctx context.Context, _ *priceClient.Request) (*priceClient.Records, error) {
	records, err := database.PriceUpdateRequests();
	response := &priceClient.Records{}

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

func (s *server) NotifyRecordsPicked(ctx context.Context, request *priceClient.Records) (*priceClient.Response, error) {
	records := request.GetEntries();
	response := &priceClient.Response{}

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

func (s *server) NotifyRecordsProcessed(ctx context.Context, request *priceClient.Records) (*priceClient.Response, error) {
	records := request.GetEntries();
	tx, err := database.GetDb().Begin();
	message := fmt.Sprintf("Successfully changed status to completed and set %v to latest", records);
	response := &priceClient.Response{Message:message}

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

func (s *server) LatestRecords(ctx context.Context, _ *priceClient.Request) (*priceClient.Records, error) {
	records, err := database.AllLatestRecords();
	response := &priceClient.Records{}

	if (err != nil) {
		log.Printf("Query failed while selecting all records \n err: %v", err)
	} else {
		for records.Next() {
			var product_id int32;
			var version string;
			var product_name string;
			var cost int32;
			var product_status string;
			var is_latest bool;
			records.Scan(&product_id, &product_name, &cost, &version, &product_status, &is_latest)
			record := priceClient.Entry{
				ProductId: product_id,
				ProductName:product_name,
				Version: version,
				Cost: cost,
				Status:product_status,
				IsLatest: is_latest}
			response.Entries = append(response.Entries, &record)
		}
	}
	return response, err
}

func (s *server) InsertRecord(ctx context.Context, request *priceClient.Record) (*priceClient.Response, error) {
	tx, err := database.GetDb().Begin();
	response := &priceClient.Response{Message:"Successfully inserted new request"}
	if (err != nil) {
		response.Message = fmt.Sprintln("Error while creating database transection");
		return response, err;
	}
	err = database.SaveEntryForUpdate(tx, request);
	if (err != nil) {
		tx.Rollback();
		response.Message = "Failed to insert new request"
		return response, err
	}
	tx.Commit();
	return response, err
}