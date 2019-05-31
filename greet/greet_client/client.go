package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/grpc/greet/greetpb"
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
	c := greetpb.NewGreetServiceClient(conn)

	//doUnary(c)
	//doServerStreaming(c)
	//doClientStreaming(c)
	doBiDiStreaming(c)
}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("starting to do bi di stream rpc")
	requests := []*greetpb.GreetEveryOneRequest{
		&greetpb.GreetEveryOneRequest{
			Greeting: "sunny",
		},
		&greetpb.GreetEveryOneRequest{
			Greeting: "sunny",
		},
		&greetpb.GreetEveryOneRequest{
			Greeting: "sunny",
		},
	}
	// we create a stream by invoking the client
	stream, err := c.GreetEveryOne(context.Background())
	if err != nil {
		log.Fatalf("err while creating stream %v", err)
		return
	}
	waitc := make(chan struct{})
	// we send some messages to clients
	go func() {
		defer stream.CloseSend()
		for _, req := range requests {
			log.Printf("sending message %v", req)
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
			fmt.Printf("recieved %v", resp.GetResult())
		}

	}()

	// block until all is done
	<-waitc

	log.Print("All closed")
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("starting to do client stream rpc")
	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "sunny",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "sunny",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "sunny",
			},
		},
	}
	stream, err := c.LongGreet(context.Background())
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

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("starting to do server stream rpc")
	req := &greetpb.GreetManyTimesRequest{
		Greeting: "Sunny",
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("err while calling greet many times rpc %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("err while reading the stream %v", err)
		}
		log.Printf("Response from Greet Many Times %v", msg.GetResult())
	}

}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("starting to do unary rpc")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "sunny",
			LastName:  "karira",
		},
	}
	resp, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("err while calling greet rpc %v", err)
	}
	fmt.Printf("resp from greet %v", resp.Result)
}
