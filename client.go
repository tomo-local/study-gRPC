package main

import (
	"example.com/myproject/chat"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := chat.NewChatServiceClient(conn)

	resp, err := c.SayHello(context.Background(), &chat.Message{
		Body: "Hello From Client!",
	})
	if err != nil {
		log.Fatalf("Error when calling SyaHello: %s", err)
	}
	log.Printf("Responce from server: %s", resp.Body)
}