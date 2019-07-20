package main

import (
	"context"
	"fmt"

	"github.com/arkiant/grpc-go-course/auth/authpb"
)

type server struct{}

func (*server) LoginUser(ctx context.Context, req *authpb.LoginUserRequest) (*authpb.LoginUserResponse, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (*server) CheckUser(ctx context.Context, req *authpb.CheckUserRequest) (*authpb.CheckUserResponse, error) {
	return nil, fmt.Errorf("Not implemented")
}
