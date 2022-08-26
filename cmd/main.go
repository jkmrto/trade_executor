package main

import (
	"fmt"

	"github.com/jkmrto/trade_executor/domain"
	"github.com/jkmrto/trade_executor/infra/binance"
)

const symbol = "BNBUSDT"

func main() {
	price := 45.0
	quantity := 45.0
	sellOrder := domain.NewSellOrder(symbol, price, quantity)
	fmt.Printf("sellOrder: %+v\n", sellOrder)

	bidUpdatesCh := make(chan []domain.Bid)
	listener := binance.BinanceListener{BidUpdatesCh: bidUpdatesCh}

	go func() { listener.Start() }()

	for bidUpdates := range listener.BidUpdatesCh {
		fmt.Printf("Proccessing the bid Updates: %+v\n", bidUpdates)
	}

	//	// use stopC to exit
	//	go func() {
	//		time.Sleep(20 * time.Second)
	//		listener.StopCh <- struct{}{}
	//	}()
	//	// remove this if you do not want to be blocked here
	//	<-listener.StopDoneCh
}
