package main

import (
	"context"
	"log"
	"time"

	"github.com/shirou/gopsutil/cpu"
	pb "path/to/generated/logdata" // Замініть на реальний шлях до logdata.pb.go

	"google.golang.org/grpc"
)

func getSystemLoad() (float64, error) {
	// Отримуємо завантаження CPU за останню секунду
	loads, err := cpu.Percent(1*time.Second, false)
	if err != nil {
		return 0, err
	}

	// Повертаємо перше значення завантаження (середнє для всіх ядер)
	return loads[0], nil
}

func main() {
	// Підключення до сервера
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Не вдалося підключитися до сервера: %v", err)
	}
	defer conn.Close()

	client := pb.NewLoggerClient(conn)

	// Отримання поточного завантаження системи
	load, err := getSystemLoad()
	if err != nil {
		log.Fatalf("Помилка під час отримання навантаження системи: %v", err)
	}

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
