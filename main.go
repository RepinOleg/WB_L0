package main

import (
	"database/sql"
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
}
