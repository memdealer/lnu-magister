package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"
)

// Генерація випадкового значення навантаження
func generateLoad() float64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Float64() * 100.0 // Випадкове значення від 0 до 100
}

// Розрахунок часу відповіді на основі навантаження
func calculateResponseTime(load float64) float64 {
	return 10 + (load * 2) // Час відповіді = базовий час + залежність від навантаження
}

// Запис нового рядка у CSV
func appendToCSV(load, responseTime float64) error {
	filePath := "server_load_data.csv"

	// Відкриваємо файл для запису, створюємо якщо не існує
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Підготовка CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Запис заголовка, якщо файл порожній
	fileInfo, _ := file.Stat()
	if fileInfo.Size() == 0 {
		writer.Write([]string{"Load", "Response_Time"})
	}

	// Запис нового рядка
	return writer.Write([]string{
		fmt.Sprintf("%.2f", load),
		fmt.Sprintf("%.2f", responseTime),
	})
}

// Обробник запиту для запису нового рядка у CSV
func handleRequest(w http.ResponseWriter, r *http.Request) {
	load := generateLoad()
	responseTime := calculateResponseTime(load)

	// Запис у файл
	err := appendToCSV(load, responseTime)
	if err != nil {
		http.Error(w, "Не вдалося записати дані у файл", http.StatusInternalServerError)
		fmt.Println("Помилка:", err)
		return
	}

	// Відповідь клієнту
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Новий рядок додано: Load=%.2f, Response_Time=%.2f\n", load, responseTime)))
	fmt.Printf("Новий рядок записано: Load=%.2f, Response_Time=%.2f\n", load, responseTime)
}

func main() {
	// HTTP маршрут
	http.HandleFunc("/log", handleRequest)

	// Запуск сервера
	fmt.Println("Сервер запущено на http://localhost:8080")
	fmt.Println("Кожен запит на /log додає новий рядок у CSV")
	http.ListenAndServe(":8080", nil)
}
