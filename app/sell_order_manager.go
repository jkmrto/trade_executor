package app

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jkmrto/trade_executor/domain"
)

// SellOrderManager ...
type SellOrderManager struct {
	ID                uuid.UUID
	SellOrder         *domain.SellOrder
	ProcessBidHandler ProcessBidHandler
	BidsCh            chan []domain.Bid
	FinishedCh        chan uuid.UUID
}

// NewSellOrderManager is a constructor
func NewSellOrderManager(so *domain.SellOrder, processBidHandler ProcessBidHandler, orderIsSoldCh chan uuid.UUID) SellOrderManager {
	return SellOrderManager{
		ID:                so.ID,
		SellOrder:         so,
		ProcessBidHandler: processBidHandler,
		BidsCh:            make(chan []domain.Bid),
		FinishedCh:        orderIsSoldCh,
	}

}

// ProcessBids ...
func (som SellOrderManager) ProcessBids() {
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
