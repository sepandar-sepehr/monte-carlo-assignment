package main

import (
	"context"
	"html/template"
	"log"
	"monte-carlo-assignment/models"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Setting timeout 30 seconds for all APIs
	ctx := context.Background()
	ctx, cancelCtx := context.WithTimeout(ctx, 30*time.Second)
	defer cancelCtx()

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.QuotePrice{})

	http.HandleFunc("/observability", renderPage)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func renderPage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("web/index.html")
	err := t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
