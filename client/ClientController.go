package client

import (
	workflow "github.com/RetailMarket/workFlowClient"
	"log"
	"google.golang.org/grpc"
)

const (
	address = "localhost:7000"
)

func CreateClientConnection() (workflow.WorkFlowClient, *grpc.ClientConn) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return workflow.NewWorkFlowClient(conn), conn
}
