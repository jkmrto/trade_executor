package main

import (
	"fmt"

	"github.com/jkmrto/trade_executor/app"
	"github.com/jkmrto/trade_executor/domain"
	"github.com/jkmrto/trade_executor/infra/binance"
)

const symbol = "BNBUSDT"

func main() {
	price := 45.0
	quantity := 45.0
	sellOrder := domain.NewSellOrder(symbol, price, quantity)
	fmt.Printf("sellOrder: %+v\n", sellOrder)

	bidsCh := make(chan []domain.Bid)
	binanceListener := binance.BinanceListener{BidsCh: bidsCh}

	go func() { binanceListener.Start() }()

	processBid := app.ProcessBidHandler{
		Exchange: app.DummyExchange{},
	}

	sellOrderManager := SellOrderManager{
		SellOrder:         &sellOrder,
		processBidHandler: processBid,
		bidsCh:            bidsCh,
	}
	sellOrderManager.processBids()

	//	// use stopC to exit
	//	go func() {
	//		time.Sleep(20 * time.Second)
	//		listener.StopCh <- struct{}{}
	//	}()
	//	// remove this if you do not want to be blocked here
	//	<-listener.StopDoneCh
}

// SellOrderManager ...
type SellOrderManager struct {
	SellOrder         *domain.SellOrder
	processBidHandler app.ProcessBidHandler
	bidsCh            chan []domain.Bid
}

func (som SellOrderManager) processBids() {
	for bids := range som.bidsCh {
		for _, bid := range bids {
			som.processBidHandler.Handle(som.SellOrder, bid)
			fmt.Printf("Sell Order: %+v\n", som.SellOrder)
			if som.SellOrder.RemainingQuantity == 0 {
				return
			}
		}
	}
}
