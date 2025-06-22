package main

import (
	"context"
	"log"
	"time"

	pb "memo/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial(
		"localhost:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewMemoServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.CreateMemo(ctx, &pb.CreateMemoRequest{
		Title:   "Test Title",
		Content: "Test Content",
	})
	if err != nil {
		log.Fatalf("could not create memo: %v", err)
	}

	log.Printf("Created Memo: %v", res.GetMemo())
}
