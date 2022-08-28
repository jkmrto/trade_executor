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

// SellOrderManager ...
type SellOrderManager struct {
	ProcessBidHandler ProcessBidHandler
	SymbolToBidRouter map[string]*BidsRouter
}

// NewSellOrderManager is a constructor
func NewSellOrderManager(processBid ProcessBidHandler, symbolToBidRouter map[string]*BidsRouter) SellOrderManager {

	return SellOrderManager{
		ProcessBidHandler: processBid,
		SymbolToBidRouter: symbolToBidRouter,
	}
}

// LaunchNewSellOrderManagerOrganizer ...
// TODO: Share context as argument for graceful exit of the OrderManager
func (somOrganizer SellOrderManager) LaunchNewSellOrderExecutor(sellOrder domain.SellOrder) error {
	bidsRouter, ok := somOrganizer.SymbolToBidRouter[sellOrder.Symbol]
	if !ok {
		return UnsupportedSymbolError{Symbol: sellOrder.Symbol}
	}

	sellOrderExecutor := NewSellOrderExecutor(&sellOrder, somOrganizer.ProcessBidHandler, bidsRouter.SoExecutorFinishedIDCh)

	go func() { sellOrderExecutor.ProcessBids() }()
	bidsRouter.NewSellOrderExecutorCh <- &sellOrderExecutor

	return nil
}
