package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	lp := filepath.Join("client", "html", "home.html")
	tmpl, err := template.ParseFiles(lp)
	if err != nil {
		log.Printf(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Printf(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Println("Template home.html executed successful!")
}
