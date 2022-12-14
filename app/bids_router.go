package app

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jkmrto/trade_executor/domain"
)

// BidsRouter allows to broadcast the bids for a symbol
// to all the sell orders opened for that symbol
type BidsRouter struct {
	BidsCh                 chan []domain.Bid
	SoExecutors            []*SellOrderExecutor
	SoExecutorFinishedIDCh chan uuid.UUID
	NewSellOrderExecutorCh chan *SellOrderExecutor

	// This attribute is only useful for debugging purposes
	Symbol string
}

// NewBidsRouter is a constructor
func NewBidsRouter(symbol string) *BidsRouter {
	return &BidsRouter{
		BidsCh:                 make(chan []domain.Bid),
		SoExecutorFinishedIDCh: make(chan uuid.UUID),
		NewSellOrderExecutorCh: make(chan *SellOrderExecutor),
		Symbol:                 symbol,
	}
}

// Start ...
func (br *BidsRouter) Start() {
	for {
		select {
		case orderExecutorPtr := <-br.NewSellOrderExecutorCh:
			br.SoExecutors = append(br.SoExecutors, orderExecutorPtr)
			fmt.Printf("[BidsRouter %+v] Added Executor: %+v, Total SellOrderExecutors: %+v\n", br.Symbol, orderExecutorPtr.ID, len(br.SoExecutors))

		case sellOrderManagerFinishedID := <-br.SoExecutorFinishedIDCh:
			index := findSellOrderManagerIndex(br.SoExecutors, sellOrderManagerFinishedID)
			br.SoExecutors = removeSellManagerAtIndex(br.SoExecutors, index)

			fmt.Printf("[BidsRouter %+v] RemainingSellOrdwerExecutors: %+v\n", br.Symbol, len(br.SoExecutors))

		case bids := <-br.BidsCh:
			for _, sellOrderManager := range br.SoExecutors {
				sellOrderManager.BidsCh <- bids
			}

		}
	}

}

func findSellOrderManagerIndex(soManagers []*SellOrderExecutor, sellOrderManagerFinishedID uuid.UUID) int {
	for index, soManager := range soManagers {
		if soManager.ID == sellOrderManagerFinishedID {
			return index
		}
	}

	// TODO: Maybe handle this with an error
	return -1
}

func removeSellManagerAtIndex(s []*SellOrderExecutor, index int) []*SellOrderExecutor {
	return append(s[:index], s[index+1:]...)
}
