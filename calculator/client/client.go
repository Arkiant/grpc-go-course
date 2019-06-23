package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

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

	//Sum(c)

	//PrimeNum(c)

	//Average(c)

	FindMaximum(c)

}

func Sum(c pb.CalculatorServiceClient) {
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

func PrimeNum(c pb.CalculatorServiceClient) {
	req := &pb.PrimeRequest{
		Num: 120.0,
	}

	res, err := c.PrimeNumber(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while call PrimeNumber function: %v", err)
	}

	for {
		msg, err := res.Recv()

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while getting data from streaming: %v", err)
		}

		log.Printf("%s ", msg.GetNum())
	}
}

func Average(c pb.CalculatorServiceClient) {
	stream, err := c.Average(context.Background())
	if err != nil {
		log.Fatalf("Error getting client: %v", err)
	}

	numbers := []int32{1, 2, 3, 4}

	for _, num := range numbers {
		stream.Send(&pb.AverageRequest{Num: num})
		time.Sleep(1 * time.Second)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving data: %v", err)
	}

	fmt.Printf("Average: %f", res.GetNum())
}

/*
FindMaximum function is a bidirectional stream implementation, this function send multiple integer numbers in parallel and receive the maximum number by the server
*/
func FindMaximum(c pb.CalculatorServiceClient) {
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Printf("Error while creating stream: %v", err)
		return
	}

	waitc := make(chan struct{})

	// Send go routine
	go func() {
		numbers := []int32{1, 5, 3, 6, 2, 20}
		for _, num := range numbers {
			stream.Send(&pb.FindMaximumRequest{
				Num: num,
			})
		}
		stream.CloseSend()
	}()

	// Receive go routine
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v", err)
				break
			}

			fmt.Printf("Received: %d\n", res.GetNum())
		}

		close(waitc)
	}()

	<-waitc
}
