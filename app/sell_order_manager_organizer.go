package app

import (
	"fmt"

	"github.com/jkmrto/trade_executor/domain"
)

// UnsupportedSymbolError is self explanatory
type UnsupportedSymbolError struct {
	Symbol string
}

func (e UnsupportedSymbolError) Error() string {
	return fmt.Sprintf("Unsupported symbol \"%s\"", e.Symbol)
}

// SellOrderManagerOrganizer ...
type SellOrderManagerOrganizer struct {
	ProcessBidHandler ProcessBidHandler
	SymbolToBidRouter map[string]*BidsRouter
}

// NewSellOrderManagerOrganizer ...
func NewSellOrderManagerOrganizer(processBid ProcessBidHandler, symbolToBidRouter map[string]*BidsRouter) SellOrderManagerOrganizer {

	return SellOrderManagerOrganizer{
		ProcessBidHandler: processBid,
		SymbolToBidRouter: symbolToBidRouter,
	}
}

// LaunchNewSellOrderManagerOrganizer ...
// TODO: Share context as argument for graceful exit of the OrderManager
func (somOrganizer SellOrderManagerOrganizer) LaunchNewSellOrderManager(sellOrder domain.SellOrder) error {
	bidsRouter, ok := somOrganizer.SymbolToBidRouter[sellOrder.Symbol]
	if !ok {
		return UnsupportedSymbolError{Symbol: sellOrder.Symbol}
	}

	sellOrderManager := NewSellOrderManager(&sellOrder, somOrganizer.ProcessBidHandler, bidsRouter.SoManagerFinishedIDCh)

	go func() { sellOrderManager.ProcessBids() }()
	bidsRouter.NewSellOrderManagerCh <- &sellOrderManager

	return nil
}
