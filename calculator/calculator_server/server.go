package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"net"

	"google.golang.org/grpc/codes"

	"github.com/grpc/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
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

func (s *server) Squareroot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {

	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("recieved neagtive number %v", number),
		)
	}
	return &calculatorpb.SquareRootResponse{
		Root: math.Sqrt(float64(number)),
	}, nil
}

func main() {

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}
	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServer(s, &server{})

	// Register reflection service on gRPC server.
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}
}
