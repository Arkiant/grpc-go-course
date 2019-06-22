package main

import (
	"context"
	"io"
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
