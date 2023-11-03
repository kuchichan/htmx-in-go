package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
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

type StockRepo struct {
	stockAssets *[]StockAsset
	nextID      int
}

type application struct {
	stockRepo *StockRepo
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

func (application *application) assets(w http.ResponseWriter, r *http.Request) {
	quantity, _ := strconv.Atoi(r.PostFormValue("quantity"))
	price, _ := strconv.Atoi(r.PostFormValue("price"))
	ticker := r.PostFormValue("ticker")
	if quantity == 0 || price == 0 || ticker == "" {
		return
	}

	asset := application.stockRepo.insert(ticker, quantity, price)

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

func main() {
	router := httprouter.New()
	app := &application{
		stockRepo: initStockAssets(),
	}

	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodPost, "/assets/", app.assets)

	server := &http.Server{
		Addr:    ":4000",
		Handler: router,
	}

	log.Fatal(server.ListenAndServe())
}

func initStockAssets() *StockRepo {
	return &StockRepo{
		stockAssets: &[]StockAsset{
			{ID: 1, Ticker: "AAPL", Price: 2990, Quantity: 100},
			{ID: 2, Ticker: "TSLA", Price: 19901, Quantity: 100},
			{ID: 3, Ticker: "JPM", Price: 19190, Quantity: 100},
		}, nextID: 4,
	}
}

func (stockRepo *StockRepo) insert(ticker string, quantity, price int) StockAsset {
	time.Sleep(1 * time.Second) // to simulate though work
	asset := StockAsset{ID: stockRepo.nextID, Ticker: ticker, Quantity: quantity, Price: price}
	*stockRepo.stockAssets = append(*stockRepo.stockAssets, asset)
	stockRepo.nextID += 1

	return asset
}
