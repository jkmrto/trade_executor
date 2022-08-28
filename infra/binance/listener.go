package binance

import (
	"fmt"

	"github.com/adshao/go-binance/v2"
	"github.com/jkmrto/trade_executor/domain"
)

// Listener consumes events fron binance WebSocket
type Listener struct {
	Symbol     string
	BidsCh     chan []domain.Bid
	StopCh     chan struct{}
	StopDoneCh chan struct{}
}

// NewListener  is a construtor
func NewListener(symbol string, bidsCh chan []domain.Bid) Listener {
	return Listener{
		BidsCh: bidsCh,
		Symbol: symbol,
	}
}

// Start begins to consume the events from binance
// It is a blocking call
func (bs Listener) Start() {
	errHandler := func(err error) {
		fmt.Println(err)
	}

	_, _, err := binance.WsDepthServe(bs.Symbol, bs.processEvent, errHandler)
	if err != nil {
		fmt.Println(err)
		return
	}

	// TODO: Handle properly the shut down of the binance listeners
	// bs.StopDoneCh = doneCh
	// bs.StopCh = stopCh

	fmt.Printf("[Binance consumer %s] Started Correctly\n", bs.Symbol)
}

func (bs Listener) processEvent(event *binance.WsDepthEvent) {

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
