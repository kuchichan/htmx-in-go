package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

type TemplateData struct {
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

type config struct {
	port int
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}

type application struct {
	stockRepo *StockRepo
	decoder   *form.Decoder
	config    config
}

func main() {
	router := httprouter.New()
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "Web app server port")
	flag.StringVar(&cfg.db.dsn,
		"db-dsn",
		"postgres://foliage:password@127.0.0.1/foliage?port=54043&sslmode=disable",
		"PostgreSQL DSN")

	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-i-conns", 25, "PostgreSQL max open connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-i-time", "15m", "PostgreSQL max connection idle time")

	app := &application{
		stockRepo: initStockAssets(),
		decoder:   form.NewDecoder(),
		config:    cfg,
	}

	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodPost, "/assets/", app.assetsCreate)

	_, err := openDB(app.config)
	if err != nil {
		fmt.Println("DB misconfigured")
		return
	}

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

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	err = db.PingContext(ctx)

	if err != nil {
		return nil, err
	}

	return db, nil
}
