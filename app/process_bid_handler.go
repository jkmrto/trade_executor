package app

import (
	"fmt"

	"github.com/jkmrto/trade_executor/domain"
)

// ProcessBidHandler is a handler for processing a new bid for a a given sell order.
type ProcessBidHandler struct {
	Exchange Exchange
}

// NewProcessBidHandler is a constructor
func NewProcessBidHandler(exchange Exchange) ProcessBidHandler {
	return ProcessBidHandler{
		Exchange: exchange,
	}
}

// Handle will do a sell of an order if the bid price is good enough
func (pbh ProcessBidHandler) Handle(sellOrder *domain.SellOrder, bid domain.Bid) error {
	sellOrderBook := sellOrder.ApplySellOnBid(bid)
	if sellOrderBook == nil {
		return nil
	}

	if err := pbh.Exchange.ApplySell(*sellOrderBook); err != nil {
		return fmt.Errorf("Error on sell order exchange: %w", err)
	}

	return nil
}
