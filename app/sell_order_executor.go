package app

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jkmrto/trade_executor/domain"
)

// SellOrderExecutor processess the incoming bid updates for selling a give order
type SellOrderExecutor struct {
	ID                          uuid.UUID
	SellOrder                   *domain.SellOrder
	ProcessBidHandler           ProcessBidHandler
	ShowSellOrderSummaryHandler ShowSellOrderSummaryHandler
	BidsCh                      chan []domain.Bid
	FinishedCh                  chan uuid.UUID
}

// NewSellOrderExecutor is a constructor
func NewSellOrderExecutor(so *domain.SellOrder, processBidHandler ProcessBidHandler, showSellOrderSummaryHandler ShowSellOrderSummaryHandler, finishedCh chan uuid.UUID) SellOrderExecutor {
	return SellOrderExecutor{
		ID:                          so.ID,
		SellOrder:                   so,
		ProcessBidHandler:           processBidHandler,
		ShowSellOrderSummaryHandler: showSellOrderSummaryHandler,
		BidsCh:                      make(chan []domain.Bid),
		FinishedCh:                  finishedCh,
	}
}

func (soe SellOrderExecutor) ProcessBids() {
	for bids := range soe.BidsCh {
		for _, bid := range bids {
			if err := soe.ProcessBidHandler.Handle(soe.SellOrder, bid); err != nil {
				fmt.Printf("[SellOrderExecutor %+v] Error processing bid: %v\n", soe.ID, err)
			}

			if soe.SellOrder.RemainingQuantity == 0 {
				soe.handleOrderWasSold()
				return
			}
		}
	}
}

func (soe SellOrderExecutor) handleOrderWasSold() {
	fmt.Printf("[SellOrderExecutor %+v] Exiting \n", soe.ID)

	if err := soe.ShowSellOrderSummaryHandler.Handle(*soe.SellOrder); err != nil {

		fmt.Printf("[SellOrderExecutor %+v] Error showing the summary for a sell oreshowing the summary for a sell orerr: %v\n", soe.ID, err)
	}
	soe.FinishedCh <- soe.ID
}
