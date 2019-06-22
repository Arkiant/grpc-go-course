package main

import (
	"context"
	"log"
	"net"

	"github.com/arkiant/grpc-go-course/calculator/pb"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *pb.NumRequest) (*pb.NumResponse, error) {
	num1 := req.GetNum1()
	num2 := req.GetNum2()

	result := num1 + num2

	res := &pb.NumResponse{
		Result: result,
	}

	return res, nil
}

func main() {

	config := pb.GetSettings()

	lis, err := net.Listen("tcp", config.Address)
	if err != nil {
		log.Fatalf("Failed to listen: %s", config.Address)
	}

	s := grpc.NewServer()

	pb.RegisterCalculatorServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
