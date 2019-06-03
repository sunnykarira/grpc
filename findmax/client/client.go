package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/grpc/findmax/findmax"
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
	c := findmax.NewFindMaxClient(conn)

	doBiDiStreaming(c)
}

func doBiDiStreaming(c findmax.FindMaxClient) {

	fmt.Println("starting to do bi di stream rpc")
	requests := []*findmax.FindMaxRequest{
		&findmax.FindMaxRequest{
			Number: 0,
		},
		&findmax.FindMaxRequest{
			Number: 1000,
		},
		&findmax.FindMaxRequest{
			Number: 15,
		},
		&findmax.FindMaxRequest{
			Number: 1001,
		},
		&findmax.FindMaxRequest{
			Number: 1003,
		},
		&findmax.FindMaxRequest{
			Number: 10000,
		},
	}

	stream, err := c.FindMax(context.Background())
	if err != nil {
		log.Fatalf("err while creating stream %v", err)
		return
	}

	waitc := make(chan struct{})
	// we send some messages to clients
	go func() {
		defer stream.CloseSend()
		for _, req := range requests {
			log.Printf("sending message %v\n", req)
			stream.Send(req)
			time.Sleep(1 * time.Second)
		}
	}()

	// we recieve a bunch of messages from server
	go func() {

		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				break
			}
			if err != nil {
				log.Fatalf("err while recieving stream %v", err)
				return
			}
			fmt.Printf("recieved %v\n", resp.GetMaximumNumber())
		}

	}()

	// block until all is done
	<-waitc

	log.Print("All closed")
}
