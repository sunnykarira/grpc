package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/grpc/computeaverage/computeaverage"
	"google.golang.org/grpc"
)

func main() {

	fmt.Println("Hello I am client")
	// grpc.WithInsecure() // no SSL
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("dial %v", err)
	}
	defer conn.Close()
	c := computeaverage.NewComputeAverageClient(conn)

	//doUnary(c)
	//doServerStreaming(c)
	doClientStreaming(c)

}

func doClientStreaming(c computeaverage.ComputeAverageClient) {
	fmt.Println("starting to do client stream rpc")
	requests := []*computeaverage.ComputeAverageRequest{
		&computeaverage.ComputeAverageRequest{
			Number: 1,
		},
		&computeaverage.ComputeAverageRequest{
			Number: 2,
		},
		&computeaverage.ComputeAverageRequest{
			Number: 3,
		},
		&computeaverage.ComputeAverageRequest{
			Number: 4,
		},
	}
	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("error while calling long greet rpc %v", err)
	}

	for _, req := range requests {
		log.Printf("Sending request %v", req)
		stream.Send(req)
		time.Sleep(1 * time.Second)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while recieving response from long greet rpc %v", err)
	}

	log.Printf("long greet reponse %v\n", res)
}
