package main

import (
	"fmt"
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

	fmt.Printf("created client %f", c)

}
