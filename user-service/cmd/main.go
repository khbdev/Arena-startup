package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"user-service/internal/config"
	"user-service/internal/handler"
	"user-service/internal/repostroy/mysql"

	"user-service/internal/usecase"

	"github.com/joho/godotenv"
	userpb "github.com/khbdev/arena-startup-proto/proto/user"
	"google.golang.org/grpc"
)

func main() {
	// 1️⃣ Load .env
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}
	config.InitRedis()
	// 2️⃣ MySQL connection
	db := config.NewMySQLConnection()
	fmt.Println("MySQL Connection Successful")

	// 3️⃣ Repository
	userRepo := mysql.NewUserRepository(db)

	// 4️⃣ Usecase
	userUsecase := usecase.NewUserUsecase(userRepo)

	// 5️⃣ Handler
	userHandler := handler.NewUserHandler(userUsecase)

	// 6️⃣ gRPC server
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, userHandler)

	fmt.Printf("gRPC server listening on port %s...\n", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
