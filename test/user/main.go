package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	userpb "github.com/khbdev/arena-startup-proto/proto/user"
	"google.golang.org/grpc"
)

func main() {
	// 1️⃣ gRPC server manzili
	serverAddr := "localhost:50051"

	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := userpb.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\nChoose action:")
		fmt.Println("1 - GetUserByTelegramId")
		fmt.Println("2 - CreateUser")
		fmt.Println("0 - Exit")
		fmt.Print("> ")

		choiceStr, _ := reader.ReadString('\n')
		choiceStr = strings.TrimSpace(choiceStr)

		switch choiceStr {
		case "1":
			// GetUserByTelegramId
			fmt.Print("Enter TelegramID: ")
			idStr, _ := reader.ReadString('\n')
			idStr = strings.TrimSpace(idStr)
			tid, _ := strconv.ParseInt(idStr, 10, 64)

			resp, err := client.GetUserByTelegramId(ctx, &userpb.GetUserRequest{
				TelegramId: tid,
			})
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			if resp.User == nil {
				fmt.Println("User not found")
			} else {
				fmt.Printf("User: %+v\n", resp.User)
			}

		case "2":
			// CreateUser
			fmt.Print("Enter TelegramID: ")
			idStr, _ := reader.ReadString('\n')
			idStr = strings.TrimSpace(idStr)
			tid, _ := strconv.ParseInt(idStr, 10, 64)

			fmt.Print("Enter Role: ")
			role, _ := reader.ReadString('\n')
			role = strings.TrimSpace(role)

			fmt.Print("Enter FirstName: ")
			firstName, _ := reader.ReadString('\n')
			firstName = strings.TrimSpace(firstName)

			resp, err := client.CreateUser(ctx, &userpb.CreateUserRequest{
				TelegramId: tid,
				Role:       role,
				FirstName:  firstName,
			})
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			fmt.Printf("Created User: %+v\n", resp.User)

		case "0":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice")
		}
	}
}
