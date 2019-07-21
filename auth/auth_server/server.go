package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/arkiant/grpc-go-course/auth/authpb"
	"github.com/arkiant/grpc-go-course/auth/database"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

func (*server) LoginUser(ctx context.Context, req *authpb.LoginUserRequest) (*authpb.LoginUserResponse, error) {

	username := req.GetUsername()
	password := req.GetPassword()

	user, err := database.LoginUser(username, password)
	if err != nil {
		return nil, fmt.Errorf("Error: %v", err)
	}

	loginUserResponse := &authpb.LoginUserResponse{
		User: &authpb.User{
			Id:   user.ID.Hex(),
			Name: user.Name,
			Role: user.Role,
		},
	}

	return loginUserResponse, nil
}

func (*server) CheckUser(ctx context.Context, req *authpb.CheckUserRequest) (*authpb.CheckUserResponse, error) {
	return nil, fmt.Errorf("Not implemented")
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Println("Blog Service Started")
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	authpb.RegisterAuthServiceServer(s, &server{})

	reflection.Register(s)

	go func() {
		fmt.Println("Starting Server...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v\n", err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch

	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("Stopping the listener")
	lis.Close()
	fmt.Println("Closing MongoDB Connection")
	fmt.Println("End of Program")
}
