package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"

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

func (*server) PrimeNumber(req *pb.PrimeRequest, res pb.CalculatorService_PrimeNumberServer) error {
	k := 2
	n := int(req.GetNum())
	for n > 1 {
		if n%k == 0 {
			res.Send(&pb.PrimeResponse{Num: float32(k)})
			n = n / k
		} else {
			k = k + 1
		}
	}

	return nil

}

func (*server) Average(stream pb.CalculatorService_AverageServer) error {

	count := 0
	numMsg := 0

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.AverageResponse{
				Num: float32(count) / float32(numMsg),
			})
		}
		if err != nil {
			log.Fatalf("Eror receiving stream data: %v", err)
		}
		count += int(res.GetNum())
		numMsg++
	}
}

/*
FindMaximum function is a bidirectional stream implementation, this function receive multiple integers from the client and send the maximum number
*/
func (*server) FindMaximum(stream pb.CalculatorService_FindMaximumServer) error {

	maxNumber := int32(0)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading number: %v", err)
			return err
		}

		number := req.GetNum()
		if number > maxNumber {
			maxNumber = number
			sendErr := stream.Send(&pb.FindMaximumResponse{Num: maxNumber})
			if sendErr != nil {
				log.Fatalf("Error while sending data to client stream: %v", err)
				return err
			}
		}

	}

}

/*
SquareRoot is a function for error handling
this RPC whill throw an exception if the sent number is negative
The error being sent is of type INVALID_ARGUMENT
*/
func (*server) SquareRoot(ctx context.Context, req *pb.SquareRootRequest) (*pb.SquareRootResponse, error) {
	fmt.Println("Received SquareRoot RPC")
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received a negative number: %d", number),
		)
	}
	return &pb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
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
