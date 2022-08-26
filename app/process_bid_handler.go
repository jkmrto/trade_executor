package app

import (
	"fmt"

	"github.com/jkmrto/trade_executor/domain"
)

// ProcessBidHandler ...
type ProcessBidHandler struct {
	Exchange Exchange
}

// Handle ...
func (pbh ProcessBidHandler) Handle(sellOrder *domain.SellOrder, bid domain.Bid) error {
	sellBook := sellOrder.ApplySellOnBid(bid)
	if sellBook == nil {
		return nil
	}

	if err := pbh.Exchange.ApplySell(*sellBook); err != nil {
		return fmt.Errorf("Error on sell order exchange: %w", err)
	}

	return nil
}
