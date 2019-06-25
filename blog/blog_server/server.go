package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/arkiant/grpc-go-course/blog/blogpb"

	"github.com/arkiant/grpc-go-course/blog/database"
	"google.golang.org/grpc"
)

type server struct{}

func main() {

	// If we crash the go code, we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Conecting to MongoDB")
	collection := database.MongoCollection()
	fmt.Printf("Conected to: %s\n", collection.Name())

	fmt.Println("Blog Service Started")
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	blogpb.RegisterBlogServiceServer(s, &server{})

	go func() {
		fmt.Println("Starting Server...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Wait for Control C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until a signal is received
	<-ch

	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("Stopping the listener")
	lis.Close()
	fmt.Println("Closing MongoDB Connection")
	database.CloseConnection()
	fmt.Println("End of Program")
}
