package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"service/server/http-server/model"

	"service/server/cache"
	"service/server/storage/postgres"
)

type AddPageData struct {
	OrderID     string
	OrderInfo   string
	ShowMessage bool
	Message     string
}

func AddOrderPage(message string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderInfo := r.URL.Query().Get("orderInfo")
		var showMessage bool
		if message != "" {
			showMessage = true
		}

		data := AddPageData{OrderInfo: orderInfo, ShowMessage: showMessage, Message: message}

		lp := filepath.Join("client", "html", "add.html")
		tmpl, err := template.ParseFiles(lp)
		if err != nil {
			log.Printf(err.Error())
			http.Error(w, fmt.Sprintf("Internal Server Error: %s", err), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			log.Printf(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		log.Println("Template add.html executed successful!")
	}
}

func AddOrder(storage *postgres.Storage, cache *cache.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		orderInfo := r.FormValue("orderInfo")

		if orderInfo == "" {
			AddOrderPage("Fields must be filled!")(w, r)
			return
		}
		var order model.Model
		if err := json.Unmarshal([]byte(orderInfo), &order); err != nil {
			AddOrderPage("Order not a model!")(w, r)
			return
		}

		// Add to storage
		err := storage.AddOrder(order)
		if err != nil {
			log.Printf(err.Error())
			AddOrderPage(err.Error())(w, r)
			return
		}
		log.Println("Order added to db successfully!")

		// Add to cache
		byt, err := json.Marshal(order)
		if err != nil {
			log.Printf(err.Error())
		} else {
			cache.SetDefault(order.OrderUID, byt)
		}

		log.Println("Order added to cache successfully!")

		AddOrderPage("Order added successfully!")(w, r)
	}
}
