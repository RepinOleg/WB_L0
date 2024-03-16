package handler

import (
	"encoding/json"
	"fmt"
	"github.com/RepinOleg/WB_L0/internal/cache"
	"github.com/RepinOleg/WB_L0/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/nats-io/stan.go"
	"github.com/xeipuuv/gojsonschema"
	"html/template"
	"net/http"
	"os"
)

type DBHandler struct {
	DB    *sqlx.DB
	Cache *cache.Cache
}

func (d *DBHandler) MsgHandler(msg *stan.Msg) {
	res := validator(msg.Data)
	if res {
		var order models.Order
		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			fmt.Println("Incorrect json file, exit")
			return
		}
		//сохраняем order in cache
		d.Cache.SetOrder(order.OrderUID, order)
		//сохраняем запись в базу
		d.AddOrder(order.OrderUID, msg.Data)
	} else {
		fmt.Println("Incorrect json file, exit")
	}
}

func (d *DBHandler) GetOrderHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("order")

	order, ok := d.GetOrderByID(idStr)
	if ok != true {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(fmt.Sprintf("User with ID = %s not found", idStr))
		return
	}
	w.WriteHeader(http.StatusOK)
	tmpl, err := template.ParseFiles("templates/order.html")
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, order)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

func (d *DBHandler) AddOrder(id string, data []byte) {
	_, err := d.DB.Exec(`INSERT INTO orders(order_id, order_data) VALUES($1, $2)`, id, data)
	if err != nil {
		fmt.Println(err)
	}
}

func (d *DBHandler) GetOrderByID(id string) (models.Order, bool) {
	return d.Cache.GetOrderByID(id)
}

func (d *DBHandler) GetAllOrders() ([]models.Order, error) {
	var orders []models.Order

	rows, err := d.DB.Query("SELECT order_id AS id, order_data AS data FROM orders")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var jsonData string
		var id string
		if err = rows.Scan(&id, &jsonData); err != nil {
			return nil, err
		}
		var order models.Order
		if err = json.Unmarshal([]byte(jsonData), &order); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func validator(input []byte) bool {
	schemaJSON, _ := os.ReadFile("schema.json")
	schema := gojsonschema.NewStringLoader(string(schemaJSON))
	in := gojsonschema.NewStringLoader(string(input))

	res, err := gojsonschema.Validate(schema, in)
	if err != nil {
		return false
	}
	return res.Valid()
}
