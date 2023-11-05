package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-playground/form/v4"
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
	decoder   *form.Decoder
}

func main() {
	router := httprouter.New()
	app := &application{
		stockRepo: initStockAssets(),
		decoder: form.NewDecoder(),
	}

	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodPost, "/assets/", app.assetsCreate)

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
