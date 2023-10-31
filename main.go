package main

import (
	"html/template"
	"log"
	"net/http"
)

type TemplateData struct {
	Hello       string
	StockAssets *[]StockAsset
}

type StockAsset struct {
	ID       int
	Ticker   string
	Quantity int // for now
	Price    int
}

func hello(w http.ResponseWriter, r *http.Request) {
	pattern := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
	}
	ts, err := template.ParseFiles(pattern...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server error", http.StatusInternalServerError)
		return
	}

	stockAssets := &[]StockAsset{
		{ID: 1, Ticker: "AAPL", Price: 2990, Quantity: 100},
		{ID: 2, Ticker: "TSLA", Price: 19901, Quantity: 100},
		{ID: 3, Ticker: "JPM", Price: 19190, Quantity: 100},
	}

	err = ts.ExecuteTemplate(w, "base", &TemplateData{Hello: "hehe", StockAssets: stockAssets})
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server error", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/hello", hello)

	server := &http.Server{
		Addr: ":4000",
	}
	log.Fatal(server.ListenAndServe())
}
