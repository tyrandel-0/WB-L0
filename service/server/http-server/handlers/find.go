package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"service/server/cache"
	"service/server/storage/postgres"
)

type FindPageData struct {
	OrderID     string
	ShowMessage bool
	Message     string
}

// GET
func FindOrderByIDPage(message string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderID := r.URL.Query().Get("orderID")
		var showMessage bool
		if message != "" {
			showMessage = true
		}
		data := FindPageData{OrderID: orderID, ShowMessage: showMessage, Message: message}

		lp := filepath.Join("client", "html", "find.html")
		tmpl, err := template.ParseFiles(lp)
		if err != nil {
			log.Printf(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			log.Printf(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		log.Println("Template find.html executed successful!")
	}
}

// POST
func FindOrderByID(storage *postgres.Storage, cache *cache.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		orderID := r.FormValue("orderID")

		if orderID == "" {
			FindOrderByIDPage("Field must be filled!")(w, r) // Повторно отображаем страницу с предупреждением
			return
		}
		var jsonB []byte
		order, ok := cache.Get(orderID)
		if ok {
			if jsonB, ok = order.([]byte); ok {
				log.Println("Get order from cache successful!")
			} else {
				log.Println("Failed to convert []byte")
				FindOrderByIDPage("Error getting data from cache!")(w, r)
			}
		} else {
			var err error
			jsonB, err = storage.GetById(orderID)
			if err != nil {
				log.Printf(err.Error())
				FindOrderByIDPage("Error getting data from database!")(w, r)
				return
			}
		}

		if jsonB == nil {
			FindOrderByIDPage("Nothing found for this id!")(w, r) // Повторно отображаем страницу с предупреждением
			return
		}

		log.Println("Find order by ID page successful!")

		OrderDetailsPage(jsonB)(w, r)
	}
}
