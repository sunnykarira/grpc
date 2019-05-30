package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/grpc/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func (s *server) Sum(ctx context.Context, req *calculatorpb.CalculatorRequest) (*calculatorpb.CalculatorResponse, error) {
	fmt.Printf("Calculator sum function was invoked %v", req)
	firstNumber := req.GetNumberOne()
	secondNumber := req.GetNumberTwo()

	return &calculatorpb.CalculatorResponse{
		CalculatedValue: firstNumber + secondNumber,
	}, nil
}

func main() {

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}
	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}
}
