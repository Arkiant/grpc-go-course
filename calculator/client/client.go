package main

import (
	"context"
	"log"

	"github.com/arkiant/grpc-go-course/calculator/pb"

	"google.golang.org/grpc"
)

func main() {

	config := pb.GetSettings()

	cc, err := grpc.Dial(config.Address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	defer cc.Close()

	c := pb.NewCalculatorServiceClient(cc)

	req := &pb.NumRequest{
		Num1: 3.0,
		Num2: 10.0,
	}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Sum RPC: %v", err)
	}

	log.Printf("Response from Sum: %v", res.Result)
}
