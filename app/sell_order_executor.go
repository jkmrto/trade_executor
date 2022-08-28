package app

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jkmrto/trade_executor/domain"
)

// SellOrderExecutor processess the incoming bid updates for selling a give order
type SellOrderExecutor struct {
	ID                uuid.UUID
	SellOrder         *domain.SellOrder
	ProcessBidHandler ProcessBidHandler
	BidsCh            chan []domain.Bid
	FinishedCh        chan uuid.UUID
}

// NewSellOrderExecutor is a constructor
func NewSellOrderExecutor(so *domain.SellOrder, processBidHandler ProcessBidHandler, orderIsSoldCh chan uuid.UUID) SellOrderExecutor {
	return SellOrderExecutor{
		ID:                so.ID,
		SellOrder:         so,
		ProcessBidHandler: processBidHandler,
		BidsCh:            make(chan []domain.Bid),
		FinishedCh:        orderIsSoldCh,
	}

}

// ProcessBids ...
func (som SellOrderExecutor) ProcessBids() {
	for bids := range som.BidsCh {
		for _, bid := range bids {
			som.ProcessBidHandler.Handle(som.SellOrder, bid)
			fmt.Printf("[%+v] After processing the bid: %+v\n", som.ID, som.SellOrder)

			if som.SellOrder.RemainingQuantity == 0 {
				fmt.Printf("[%+v] Consumer Exiting \n", som.ID)
				som.FinishedCh <- som.ID
				return
			}
		}
	}
}
