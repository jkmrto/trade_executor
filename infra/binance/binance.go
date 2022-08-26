package binance

import (
	"fmt"

	"github.com/adshao/go-binance/v2"
	"github.com/jkmrto/trade_executor/domain"
)

const symbol = "BNBUSDT"

// BinanceListener ...
type BinanceListener struct {
	BidsCh     chan []domain.Bid
	StopCh     chan struct{}
	StopDoneCh chan struct{}
}

// Start ...
func (bs BinanceListener) Start() {
	errHandler := func(err error) {
		fmt.Println(err)
	}

	doneCh, stopCh, err := binance.WsDepthServe(symbol, bs.processBinanceEvent, errHandler)
	if err != nil {
		fmt.Println(err)
		return
	}

	bs.StopDoneCh = doneCh
	bs.StopCh = stopCh

	return
}

func (bs BinanceListener) processBinanceEvent(event *binance.WsDepthEvent) {

	var bids []domain.Bid

	// fmt.Printf("Event: %+v\n", event)
	for index, rawBid := range event.Bids {
		bidID := binanceBidID(event, index)

		bid, err := domain.NewBidFromRaw(bidID, event.Symbol, rawBid.Price, rawBid.Quantity)
		if err != nil {
			fmt.Printf("Error () processing bid %+v, %v", err, rawBid)
		}
		bids = append(bids, bid)
	}

	// Using a gorouting we avoid blocking the Binance consumer
	go func() { bs.BidsCh <- bids }()
}

func binanceBidID(event *binance.WsDepthEvent, index int) string {
	return fmt.Sprintf("%d-%d-%d", event.FirstUpdateID, event.LastUpdateID, index)
}
