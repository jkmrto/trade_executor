package main

import (
	"fmt"
	"net/http"

	"github.com/jkmrto/trade_executor/app"
	"github.com/jkmrto/trade_executor/config"
	"github.com/jkmrto/trade_executor/infra/binance"
	"github.com/jkmrto/trade_executor/infra/sqlite3"

	httpx "github.com/jkmrto/trade_executor/infra/http"
)

const symbol = "BNBUSDT"

func main() {
	conf := config.New()

	bidsRouter := app.NewBidsRouter()
	go func() { bidsRouter.Start() }()

	binanceListener := binance.BinanceListener{BidsCh: bidsRouter.BidsCh}
	go func() { binanceListener.Start() }()

	db, err := setupDB(conf.Sqlite3)
	if err != nil {
		fmt.Printf("%+v", err)
	}

	processBidHandler := app.NewProcessBidHandler(db)
	somOrganizer := app.NewSellOrderManagerOrganizer(processBidHandler, bidsRouter)

	startHTTPServer(somOrganizer, conf.HTTP)
}

// startHTTPServer is a blocking call
func startHTTPServer(somOrganizer app.SellOrderManagerOrganizer, conf httpx.Config) {
	handlers := []httpx.EndpointHandlerMethod{
		{
			Endpoint:    "/SellOrder",
			Method:      http.MethodPost,
			HandlerFunc: httpx.CreateSellOrder(somOrganizer),
		},
	}

	httpx.ListenAndServe(conf, handlers)
}

func setupDB(conf sqlite3.Config) (sqlite3.Database, error) {
	db, err := sqlite3.NewDatabase(conf)
	if err != nil {
		return sqlite3.Database{}, fmt.Errorf("error when connecting to the DB: %v", err)
	}

	//TODO This is not the ideal way of handling this error (since we are comparing two strings)
	err = db.RunMigrations()
	if err != nil && err.Error() != "no change" {
		return sqlite3.Database{}, fmt.Errorf("error running migrations: %+v", err)
	}

	return db, nil
}
