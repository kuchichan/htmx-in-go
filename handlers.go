package main

import (
	"html/template"
	"log"
	"net/http"
)

type AssetCreate struct {
	Ticker   string `form:"ticker"`
	Quantity int    `form:"quantity"`
	Price    int    `form:"price"`
}

func (application *application) home(w http.ResponseWriter, r *http.Request) {
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

	err = ts.ExecuteTemplate(w, "base", &TemplateData{Hello: "hehe", StockAssets: application.stockRepo.stockAssets})
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server error", http.StatusInternalServerError)
	}
}

func (application *application) assetsCreate(w http.ResponseWriter, r *http.Request) {
	var form AssetCreate

	err := application.decodeForm(r, &form)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	if !isNotEmpty(form.Ticker) || !isPositive(form.Quantity) || !isPositive(form.Price) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	asset := application.stockRepo.insert(form.Ticker, form.Quantity, form.Price)

	ts, err := template.ParseFiles("./ui/html/partials/asset-row.tmpl.html")
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "asset-row", &asset)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server error", http.StatusInternalServerError)
	}
}
