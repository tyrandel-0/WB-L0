package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"service/server/http-server/model"
)

func OrderDetailsPage(jsonB []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data model.Model

		json.NewDecoder(strings.NewReader(string(jsonB))).Decode(&data)

		lp := filepath.Join("client", "html", "orderDetails.html")
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

		log.Println("Template order.html executed successful!")
	}
}
