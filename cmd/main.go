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

// TODO: Handle the binance consumers exits properly
// TODO: Share context from main to all downstream components
// TODO: Maybe unblock the sending of messages on the bidsRouter
// from the bidsrouter to the sellOrders executers?

func main() {

	conf := config.New()

	db, err := setupDB(conf.Sqlite3)
	if err != nil {
		fmt.Printf("Error on setupDB: %+v\n", err)
	}

	symbolToBidsRouter := make(map[string]*app.BidsRouter)
	for _, symbol := range conf.SupportedSymbols {
		symbolToBidsRouter[symbol] = startPipelineForSymbol(symbol)
	}

	processBidHandler := app.NewProcessBidHandler(db)
	sellOrderManager := app.NewSellOrderManager(processBidHandler, symbolToBidsRouter)

	if err = startHTTPServer(sellOrderManager, conf.HTTP); err != nil {
		fmt.Printf("Error starting HTTP server: %+v\n", err)
	}
}

func startPipelineForSymbol(symbol string) *app.BidsRouter {
	bidsRouter := app.NewBidsRouter(symbol)
	go func() { bidsRouter.Start() }()

	binanceListener := binance.NewListener(symbol, bidsRouter.BidsCh)
	go func() { binanceListener.Start() }()

	return bidsRouter
}

// startHTTPServer is a blocking call
func startHTTPServer(sellOrderManager app.SellOrderManager, conf httpx.Config) error {
	handlers := []httpx.EndpointHandlerMethod{
		{
			Endpoint:    "/SellOrder",
			Method:      http.MethodPost,
			HandlerFunc: httpx.CreateSellOrder(sellOrderManager),
		},
	}

	return httpx.ListenAndServe(conf, handlers)
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
