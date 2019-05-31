package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/grpc/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct{}

func (s *server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked %v", req)
	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()

	result := "Hello " + firstName + " " + lastName
	return &greetpb.GreetResponse{
		Result: result,
	}, nil
}

func (s *server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes function was invoked %v", req)
	firstName := req.GetGreeting()
	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + " number " + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(1 * time.Second)
	}
	return nil
}

func (s *server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Printf("LongGreet was invoked")
	result := "Hello "
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// We have finished reading client stream
			return stream.SendAndClose(
				&greetpb.LongGreetResponse{
					Result: result,
				},
			)
		}
		if err != nil {
			log.Fatalf("Error while reading client stream %v", err)
		}
		firstName := req.GetGreeting().GetFirstName()
		result += firstName + " ! "
	}

}

func (s *server) GreetEveryOne(stream greetpb.GreetService_GreetEveryOneServer) error {
	fmt.Printf("GreetEveryOne was invoked")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// stop recieving stream
			return nil
		}
		if err != nil {
			log.Fatalf("err while reading client stream %v", err)
			return err
		}
		firstName := req.GetGreeting()
		result := "Hello " + firstName + " "
		err = stream.Send(&greetpb.GreetEveryOneResponse{
			Result: result,
		})
		if err != nil {
			log.Fatalf("err while sending data to client stream %v", err)
			return err
		}
	}
}

func main() {

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}
	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}
}
