package app

import (
	"fmt"

	"github.com/jkmrto/trade_executor/domain"
)

// ProcessBidHandler ...
type ProcessBidHandler struct {
	Exchange Exchange
}

// This interface is only here for testing purpouses
// It allows decouple the testing of the infra layer from the app one
// type ProcessBidInterface struct {
// 	Exchange Exchange
// }

// Handle ...
func (pbh ProcessBidHandler) Handle(sellOrder *domain.SellOrder, bid domain.Bid) error {
	sellBook := sellOrder.ApplySellOnBid(bid)
	if sellBook == nil {
		fmt.Printf("hello\n")
		return nil
	}

	if err := pbh.Exchange.ApplySell(sellBook); err != nil {
		return fmt.Errorf("Error on sell order exchange: %w", err)
	}

	return nil
}
