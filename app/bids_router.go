package app

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jkmrto/trade_executor/domain"
)

// BidsRouter ...
type BidsRouter struct {
	BidsCh                chan []domain.Bid
	SoManagers            []*SellOrderManager
	SoManagerFinishedIDCh chan uuid.UUID
	NewSellOrderManagerCh chan *SellOrderManager
}

// NewBidsRouter ...
func NewBidsRouter() BidsRouter {
	return BidsRouter{
		BidsCh:                make(chan []domain.Bid),
		SoManagerFinishedIDCh: make(chan uuid.UUID),
		NewSellOrderManagerCh: make(chan *SellOrderManager),
	}
}

// Start ...
func (br BidsRouter) Start() {
	for {
		select {
		case orderManagerPtr := <-br.NewSellOrderManagerCh:
			fmt.Printf("[BidsRouter][New Sell Order manager] %+v \n", orderManagerPtr.ID)
			br.SoManagers = append(br.SoManagers, orderManagerPtr)
		case bids := <-br.BidsCh:
			fmt.Printf("[BidsRouter][Bids received]: Sell Order Managers Actived %+v \n", len(br.SoManagers))
			for _, sellOrderManager := range br.SoManagers {
				fmt.Printf("[BidsRouter][Sending Bids]: %+v \n", sellOrderManager.ID)
				sellOrderManager.BidsCh <- bids
			}
		case sellOrderManagerFinishedID := <-br.SoManagerFinishedIDCh:
			index := findSellOrderManagerIndex(br.SoManagers, sellOrderManagerFinishedID)
			br.SoManagers = removeSellManagerAtIndex(br.SoManagers, index)
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
