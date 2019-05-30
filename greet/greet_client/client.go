package main

import (
	"context"
	"fmt"
	"io"
	"log"

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
	doServerStreaming(c)

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
