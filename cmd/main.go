package main

import (
	"fmt"
	"time"

	"github.com/jkmrto/trade_executor/app"

	"github.com/jkmrto/trade_executor/domain"
	"github.com/jkmrto/trade_executor/infra/binance"
)

const symbol = "BNBUSDT"

func main() {
	price := 45.0
	quantity := 45.0
	sellOrder1 := domain.NewSellOrder(symbol, price, quantity)
	fmt.Printf("sellOrder1: %+v\n", sellOrder1)

	price = 45.0
	quantity = 1000.0
	sellOrder2 := domain.NewSellOrder(symbol, price, quantity)
	fmt.Printf("sellOrder2: %+v\n", sellOrder2)

	price = 280.0
	quantity = 100.0
	sellOrder3 := domain.NewSellOrder(symbol, price, quantity)
	fmt.Printf("sellOrder3: %+v\n", sellOrder3)

	// list of channels? iteract over them

	bidsRouter := app.NewBidsRouter()
	go func() { bidsRouter.Start() }()

	binanceListener := binance.BinanceListener{BidsCh: bidsRouter.BidsCh}
	go func() { binanceListener.Start() }()

	processBid := app.ProcessBidHandler{
		Exchange: app.DummyExchange{},
	}

	sellOrderManager1 := app.NewSellOrderManager(&sellOrder1, processBid, bidsRouter.SoManagerFinishedIDCh)
	go func() { sellOrderManager1.ProcessBids() }()
	bidsRouter.NewSellOrderManagerCh <- &sellOrderManager1

	sellOrderManager2 := app.NewSellOrderManager(&sellOrder2, processBid, bidsRouter.SoManagerFinishedIDCh)
	go func() { sellOrderManager2.ProcessBids() }()
	bidsRouter.NewSellOrderManagerCh <- &sellOrderManager2

	sellOrderManager3 := app.NewSellOrderManager(&sellOrder3, processBid, bidsRouter.SoManagerFinishedIDCh)
	go func() { sellOrderManager3.ProcessBids() }()
	bidsRouter.NewSellOrderManagerCh <- &sellOrderManager3

	// The new SellOrderManager will arrive asynchrounously to the bids router
	// Since they weill be created form an HTTP interface

	//	// use stopC to exit
	//	go func() {                     	//   time.Sleep(20 * time.Second)
	//	}()
	//		listener.StopCh <- struct{}{}	//	// remove this if you do not want to be blocked here
	//	<-listener.StopDoneCh

	time.Sleep(time.Hour)
}
