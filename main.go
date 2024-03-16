package main

import (
	"fmt"
	"github.com/RepinOleg/WB_L0/internal/cache"
	"github.com/RepinOleg/WB_L0/internal/dbs"
	"github.com/RepinOleg/WB_L0/internal/handler"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"log"
	"net/http"
)

func main() {
	// Подключение к серверу NATS Streaming
	sc, err := stan.Connect("test-cluster", "sub")
	if err != nil {
		log.Fatalf("Error connecting to NATS Streaming: %v", err)
	}
	defer sc.Close()

	cfg := dbs.Config{
		Addr:     "localhost",
		Port:     5432,
		User:     "admin",
		Password: "admin",
		DB:       "postgres",
	}
	db, err := dbs.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	memCache := cache.NewCache()
	app := &handler.DBHandler{DB: db, Cache: memCache}

	// получаем из базы все записи
	orders, err := app.GetAllOrders()
	if err != nil {
		fmt.Println(err.Error())
	}
	//Добавляем все полученные записи в кэш
	app.Cache.SetAllOrders(orders)

	_, err = sc.Subscribe("order", app.MsgHandler, stan.DurableName("order"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	http.HandleFunc("/", app.GetOrderHandler)
	fmt.Println("starting server at :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}
