package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/grpc/computeaverage/computeaverage"
	"google.golang.org/grpc"
)

type server struct{}

func (s *server) ComputeAverage(stream computeaverage.ComputeAverage_ComputeAverageServer) error {
	fmt.Printf("Compute Average was invoked")
	res := float64(0)
	count := int64(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			if count != 0 {
				// We have finished reading client stream
				return stream.SendAndClose(
					&computeaverage.ComputeAverageResponse{
						Response: res / float64(count),
					},
				)
			}
			return stream.SendAndClose(
				&computeaverage.ComputeAverageResponse{
					Response: 0,
				},
			)
		}
		if err != nil {
			log.Fatalf("Error while reading client stream %v", err)
		}
		count++
		res += req.GetNumber()
	}
}

func main() {

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}
	s := grpc.NewServer()
	computeaverage.RegisterComputeAverageServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}
}
