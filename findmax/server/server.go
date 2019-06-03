package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/grpc/findmax/findmax"
	"google.golang.org/grpc"
)

type server struct{}

func (s *server) FindMax(stream findmax.FindMax_FindMaxServer) error {
	fmt.Printf("FindMax was invoked")
	maxNumber := int64(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// We have finished reading client stream
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream %v", err)
		}
		number := req.GetNumber()
		if number > maxNumber {
			maxNumber = number
			err = stream.Send(
				&findmax.FindMaxResponse{
					MaximumNumber: maxNumber,
				},
			)
			if err != nil {
				log.Fatalf("Error while sending data to client: %v", err)
				return err
			}
		}
	}
}

func main() {

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}
	s := grpc.NewServer()
	findmax.RegisterFindMaxServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}
}
