package service

import (
	"context"
	"fmt"

	hellopb "greeting/grpc"
)

type HelloService struct {
	hellopb.UnimplementedGreetingServiceServer
}

func (s *HelloService) Hello(ctx context.Context, req *hellopb.HelloRequest) (*hellopb.HelloResponse, error) {
	return &hellopb.HelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.GetName()),
	}, nil
}
