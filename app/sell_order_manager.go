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
	ProcessBidHandler           ProcessBidHandler
	ShowSellOrderSummaryHandler ShowSellOrderSummaryHandler
	SymbolToBidRouter           map[string]*BidsRouter
}

// NewSellOrderManager is a constructo
func NewSellOrderManager(processBid ProcessBidHandler, showSellOrderSummary ShowSellOrderSummaryHandler, symbolToBidRouter map[string]*BidsRouter) SellOrderManager {

	return SellOrderManager{
		ProcessBidHandler:           processBid,
		ShowSellOrderSummaryHandler: showSellOrderSummary,
		SymbolToBidRouter:           symbolToBidRouter,
	}
}

//  LaunchNewSellOrderExecutor ...
// TODO: Share context as argument for graceful exit of the OrderManager
func (som SellOrderManager) LaunchNewSellOrderExecutor(sellOrder domain.SellOrder) error {
	bidsRouter, ok := som.SymbolToBidRouter[sellOrder.Symbol]
	if !ok {
		return UnsupportedSymbolError{Symbol: sellOrder.Symbol}
	}

	sellOrderExecutor := NewSellOrderExecutor(&sellOrder, som.ProcessBidHandler, som.ShowSellOrderSummaryHandler, bidsRouter.SoExecutorFinishedIDCh)

	go func() { sellOrderExecutor.ProcessBids() }()
	bidsRouter.NewSellOrderExecutorCh <- &sellOrderExecutor

	return nil
}
