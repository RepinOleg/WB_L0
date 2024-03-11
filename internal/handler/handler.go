package handler

import (
	"encoding/json"
	"github.com/RepinOleg/WB_L0/internal/repository"
	"github.com/nats-io/stan.go"
	"log"
	"net/http"
	"strconv"
)

type HandleFunc func(w http.ResponseWriter, r *http.Request)

func (f HandleFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(w, r)
}
func MsgHandler(msg *stan.Msg, r *repository.Repository) {
	// Todo FIX ID
	err := r.AddOrder(1, msg.Data)
	if err != nil {
		log.Println(err, "Put Order")
	}
}

func GetOrderHandler(repo *repository.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("order")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		order, err := repo.GetOrderByID(id)
		if err != nil {
			json.NewEncoder(w).Encode(err.Error())
		} else {
			json.NewEncoder(w).Encode(order)
		}
	}
}
