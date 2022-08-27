package app

import "github.com/jkmrto/trade_executor/domain"

// SellOrderManagerOrganizer ...
type SellOrderManagerOrganizer struct {
	ProcessBidHandler ProcessBidHandler
	BidsRouter        *BidsRouter
}

// NewSellOrderManagerOrganizer ...
func NewSellOrderManagerOrganizer(processBid ProcessBidHandler, bidsRouter *BidsRouter) SellOrderManagerOrganizer {

	return SellOrderManagerOrganizer{
		ProcessBidHandler: processBid,
		BidsRouter:        bidsRouter,
	}
}

// LaunchNewSellOrderManager ...
// TODO: Share context as argument for graceful exit of the OrderManager
func (somOrganizer SellOrderManagerOrganizer) LaunchNewSellOrderManager(sellOrder domain.SellOrder) {
	sellOrderManager := NewSellOrderManager(&sellOrder, somOrganizer.ProcessBidHandler, somOrganizer.BidsRouter.SoManagerFinishedIDCh)

	go func() { sellOrderManager.ProcessBids() }()
	somOrganizer.BidsRouter.NewSellOrderManagerCh <- &sellOrderManager

}
