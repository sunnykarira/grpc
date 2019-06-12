package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	blog "github.com/grpc/blog/blogpb"
	"google.golang.org/grpc"
)

type server struct{}

func main() {
	// If we crash the go code we get the filename name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Blog Service Started")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	blog.RegisterBlogServiceServer(s, &server{})

	go func() {
		fmt.Println("Starting the server")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	// Wait for control+C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until the signal is recieved
	<-ch
	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("Closing the listener")
	lis.Close()
	fmt.Println("everything stopped sucessfully")
}
