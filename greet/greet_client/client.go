package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/grpc/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func main() {

	fmt.Println("Hello I'm a client")

	tls := false
	opts := grpc.WithInsecure()
	if tls {
		certFile := "ssl/ca.crt" // Certificate Authority Trust certificate
		creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")
		if sslErr != nil {
			log.Fatalf("Error while loading CA trust certificate: %v", sslErr)
			return
		}
		opts = grpc.WithTransportCredentials(creds)
	}

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	// fmt.Printf("Created client: %f", c)

	doUnary(c)
	// doServerStreaming(c)
	// doClientStreaming(c)
	// doBiDiStreaming(c)

	// doUnaryWithDeadline(c, 5*time.Second) // should complete
	// doUnaryWithDeadline(c, 1*time.Second) // should timeout
}

func doUnaryWithDeadLine(c greetpb.GreetServiceClient, seconds time.Duration) {

	fmt.Println("starting to do unary with deadline rpc")
	req := &greetpb.GreetWithDeadLineRequest{
		Greeting: "sunny karira",
	}
	ctx, cancel := context.WithTimeout(context.Background(), seconds)
	defer cancel()
	resp, err := c.GreetWithDeadLine(ctx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("timeout deadline was exceeded")
			} else {
				fmt.Println("unexpected error")
			}

		} else {
			log.Fatalf("err while calling greet rpc %v", err)
		}

		return
	}
	fmt.Printf("resp from greet %v", resp.Result)

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
