package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"log"
)

func main() {
	// Подключение к серверу NATS Streaming
	sc, err := stan.Connect("test-cluster", "sub")
	if err != nil {
		log.Fatalf("Error connecting to NATS Streaming: %v", err)
	}
	defer sc.Close()
	fmt.Println("Connected to nats streaming")

	connStr := "user=admin password=admin dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to db")
	defer db.Close()
	sub, err := sc.Subscribe("order", func(msg *stan.Msg) {
		handleJson(msg, db)
	})
	defer sub.Unsubscribe()

	select {} // Бесконечный цикл для ожидания сообщений
}

func handleJson(msg *stan.Msg, db *sql.DB) {
	var jsonData map[string]interface{}
	err := json.Unmarshal(msg.Data, &jsonData)
	if err != nil {
		log.Println("Error unmarshalling JSON data:", err)
		return
	}

	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		log.Println("Error marshalling JSON data:", err)
		return
	}

	_, err = db.Exec("INSERT INTO orders (order_data) VALUES ($1)", jsonBytes)
	if err != nil {
		log.Println("Error inserting data into PostgreSQL:", err)
		return
	}

	log.Println("JSON data saved to PostgreSQL")
}
