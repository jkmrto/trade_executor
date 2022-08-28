package app

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jkmrto/trade_executor/domain"
)

// BidsRouter allows to broadcast the bids for a symbol
// to all the sell orders opened for that symbol
type BidsRouter struct {
	BidsCh                chan []domain.Bid
	SoManagers            []*SellOrderManager
	SoManagerFinishedIDCh chan uuid.UUID
	NewSellOrderManagerCh chan *SellOrderManager

	// This attribute is only useful for debugging purposes
	Symbol string
}

// NewBidsRouter is a constructor
func NewBidsRouter(symbol string) *BidsRouter {
	return &BidsRouter{
		BidsCh:                make(chan []domain.Bid),
		SoManagerFinishedIDCh: make(chan uuid.UUID),
		NewSellOrderManagerCh: make(chan *SellOrderManager),
		Symbol:                symbol,
	}
}

// Start ...
func (br *BidsRouter) Start() {
	for {
		select {
		case orderManagerPtr := <-br.NewSellOrderManagerCh:
			fmt.Printf("[BidsRouter %+v][New Sell Order manager] %+v \n", br.Symbol, orderManagerPtr.ID)
			br.SoManagers = append(br.SoManagers, orderManagerPtr)
			fmt.Printf("[BidsRouter %+v] soManagers: %+v\n", br.Symbol, len(br.SoManagers))

		case sellOrderManagerFinishedID := <-br.SoManagerFinishedIDCh:
			fmt.Printf("[BidsRouter %+v] Remove sell order manager: %+v\n", br.Symbol, sellOrderManagerFinishedID)
			index := findSellOrderManagerIndex(br.SoManagers, sellOrderManagerFinishedID)
			br.SoManagers = removeSellManagerAtIndex(br.SoManagers, index)

			fmt.Printf("[BidsRouter %+v] soManagers: %+v\n", br.Symbol, br.SoManagers)

		case bids := <-br.BidsCh:
			//		fmt.Printf("[BidsRouter][Bids received]: Sell Order Managers Actived %+v \n", len(br.SoManagers))
			for _, sellOrderManager := range br.SoManagers {
				fmt.Printf("[BidsRouter %+v]: %+v \n", br.Symbol, sellOrderManager.ID)
				sellOrderManager.BidsCh <- bids
			}

		}
	}

}

func findSellOrderManagerIndex(soManagers []*SellOrderManager, sellOrderManagerFinishedID uuid.UUID) int {
	for index, soManager := range soManagers {
		if soManager.ID == sellOrderManagerFinishedID {
			return index
		}
	}

	// TODO: Maybe handle this with an error
	return -1
}

func removeSellManagerAtIndex(s []*SellOrderManager, index int) []*SellOrderManager {
	return append(s[:index], s[index+1:]...)
}
