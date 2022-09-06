package main

// GRPC client implementation of TodoServiceClient

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	pb "github.com/yuhuishi-convect/grpc-todo/gen/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// CLI flag for grpc server host and port
var (
	host = flag.String("host", "localhost", "The server host")
	port = flag.Int("port", 8080, "The server port")
)

func main() {
	flag.Parse()

	// connect to the grpc server
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", *host, *port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// create a client
	client := pb.NewTodoServiceClient(conn)
	// create a context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// contact the server and print out its response.
	r, err := client.List(ctx, &pb.ListRequest{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	} else {
		log.Printf("Todo items: %v", r)
	}

	// create a new todo item
	createResponse, err := client.Create(ctx, &pb.CreateRequest{
		Title: "New Todo Item",
	})

	if err != nil {
		log.Fatalf("could not create todo item: %v", err)
	} else {
		log.Printf("Created todo item: %v", createResponse)
	}

	// fetch the updated todo list
	r, err = client.List(ctx, &pb.ListRequest{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	} else {
		log.Printf("Todo items: %v", r)
	}

}
