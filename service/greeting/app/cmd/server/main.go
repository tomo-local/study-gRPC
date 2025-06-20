package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	hellopb "greeting/grpc"
	"greeting/service"
)

func NewHelloService() *service.HelloService {
	return &service.HelloService{}
}

func main() {
	// 1. 8080番portのListenerを作成
	port := 8080
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	// 2. gRPCサーバーを作成
	s := grpc.NewServer()

	// 2-1. 作成したgRPCサーバーに、サービスを登録する
	hellopb.RegisterGreetingServiceServer(s, NewHelloService())

	// 2-2. gRPCサーバーに、Reflectionを登録する
	reflection.Register(s)

	// 3. 作成したgRPCサーバーを、8080番ポートで稼働させる
	go func() {
		log.Printf("start gRPC server port: %v", port)
		s.Serve(listener)
	}()

	// 4.Ctrl+Cが入力されたらGraceful shutdownされるようにする
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping gRPC server...")
	s.GracefulStop()
}
