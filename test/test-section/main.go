package main

import (
	"context"
	"fmt"
	"log"
	"time"

	testpb "github.com/khbdev/arena-startup-proto/proto/test-section"
	"google.golang.org/grpc"
)

func main() {

	
	address := "localhost:50052"


	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Ulanishda xato: %v", err)
	}
	defer conn.Close()

	client := testpb.NewResultServiceClient(conn)

	
	req := &testpb.GetUserTestResultRequest{
		TelegramId: 3555245,
		TestId:     "TST-GVt05p",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()


	resp, err := client.GetUserTestResult(ctx, req)
	if err != nil {
		log.Fatalf("RPC xato: %v", err)
	}

	
	fmt.Println("Serverdan qaytgan JSON:")
	fmt.Println(resp.JsonData)
}
