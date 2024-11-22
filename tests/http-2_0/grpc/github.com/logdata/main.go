package main

import (
	"context"
	"log"
	"time"

	pb "path/to/generated/logdata" // Замініть на реальний шлях до logdata.pb.go

	"google.golang.org/grpc"
)

func main() {
	// Підключення до сервера
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Не вдалося підключитися до сервера: %v", err)
	}
	defer conn.Close()

	client := pb.NewLoggerClient(conn)

	// Випадкове значення навантаження
	load := float64(10 + rand.Intn(90)) // Навантаження у діапазоні [10, 100]

	// Виклик RPC LogData
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.LogData(ctx, &pb.LogRequest{Load: load})
	if err != nil {
		log.Fatalf("Помилка під час виклику LogData: %v", err)
	}

	// Виведення відповіді
	log.Printf("Відповідь сервера: %s", res.GetMessage())
}
