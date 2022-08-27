package main

import (
	"net/http"

	"github.com/jkmrto/trade_executor/app"
	"github.com/jkmrto/trade_executor/config"

	"github.com/jkmrto/trade_executor/infra/binance"
	httpx "github.com/jkmrto/trade_executor/infra/http"
)

const symbol = "BNBUSDT"

func main() {
	bidsRouter := app.NewBidsRouter()
	go func() { bidsRouter.Start() }()

	binanceListener := binance.BinanceListener{BidsCh: bidsRouter.BidsCh}
	go func() { binanceListener.Start() }()

	processBid := app.ProcessBidHandler{
		Exchange: app.DummyExchange{},
	}

	somOrganizer := app.NewSellOrderManagerOrganizer(processBid, bidsRouter)

	// TODO: Handle the binance exit properly
	// The new SellOrderManager will arrive asynchrounously to the bids router
	// Since they weill be created form an HTTP interface

	//	// use stopC to exit
	//	go func() {                     	//   time.Sleep(20 * time.Second)
	//	}()
	//		listener.StopCh <- struct{}{}	//	// remove this if you do not want to be blocked here
	//	<-listener.StopDoneCh

	handlers := []httpx.EndpointHandlerMethod{
		{
			Endpoint:    "/SellOrder",
			Method:      http.MethodPost,
			HandlerFunc: httpx.CreateSellOrder(somOrganizer),
		},
	}

	conf := config.New()
	httpx.ListenAndServe(conf.HTTP, handlers)

}
