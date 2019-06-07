package main

import (
	"context"
	"fmt"
	"log"

	"github.com/grpc/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {

	fmt.Println("Hello I am client")
	// grpc.WithInsecure() // no SSL
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("dial %v", err)
	}
	defer conn.Close()
	c := calculatorpb.NewCalculatorClient(conn)

	//doUnary(c)
	doErrorUnary(c)

}

func doErrorUnary(c calculatorpb.CalculatorClient) {
	fmt.Println("starting to do error unary rpc")
	// correct call
	req := &calculatorpb.SquareRootRequest{
		Number: 1,
	}
	resp, err := c.Squareroot(context.Background(), req)
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			log.Fatalf("big error while calling square root rpc %v", respErr)
			fmt.Println(respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				log.Fatalf("invalid argument %v", err)
			}
		}
		log.Fatalf("err while calling square root rpc %v", err)
	}
	fmt.Printf("resp from square root %v", resp.Root)

	// error call
	req = &calculatorpb.SquareRootRequest{
		Number: -100,
	}
	resp, err = c.Squareroot(context.Background(), req)
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			log.Fatalf("big error while calling square root rpc %v", respErr)
			fmt.Println(respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				log.Fatalf("invalid argument %v", err)
			}
		}
		log.Fatalf("err while calling square root rpc %v", err)
	}
	fmt.Printf("resp from square root %v", resp.Root)

}

func doUnary(c calculatorpb.CalculatorClient) {
	fmt.Println("starting to do unary rpc")
	req := &calculatorpb.CalculatorRequest{
		NumberOne: 1,
		NumberTwo: 2,
	}
	resp, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("err while calling greet rpc %v", err)
	}
	fmt.Printf("resp from greet %v", resp.CalculatedValue)
}
