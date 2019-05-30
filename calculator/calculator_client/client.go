package main

import (
	"context"
	"fmt"
	"log"

	"github.com/grpc/calculator/calculatorpb"
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
	c := calculatorpb.NewCalculatorClient(conn)

	doUnary(c)

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
