package main

import (
	"fmt"
	"log"
	"net"

	"test-section-serve/internal/config"
	"test-section-serve/internal/handler"
	"test-section-serve/internal/service"

	testpb "github.com/khbdev/arena-startup-proto/proto/test-section"
	"google.golang.org/grpc"
)

func main() {
	
	config.InitRedis()


	port := config.GetEnv("GRPC_PORT", "50052")
	address := fmt.Sprintf(":%s", port)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Portni tinglashda xato: %v", err)
	}

	serv := service.NewTestService()

	
	resultHandler := handler.NewResultHandler(serv)

	
	grpcServer := grpc.NewServer()

	testpb.RegisterResultServiceServer(grpcServer, resultHandler)

	fmt.Println("gRPC server listening on port", port)

	
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC server ishlamay qoldi: %v", err)
	}
}
