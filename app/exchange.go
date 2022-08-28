package app

import (
	"fmt"

	"github.com/jkmrto/trade_executor/domain"
)

// Exchange define the contract for selling trades
type Exchange interface {
	ApplySell(domain.SellOrderBook) error
}

// DummyExchange is a dummy implenetaion for operations in an exchange
type DummyExchange struct{}

// ApplySell just prints a given sell order
func (DummyExchange) ApplySell(sb domain.SellOrderBook) error {
	fmt.Printf("[%+v] SellBook: %+v \n", sb.BidID, sb)
	return nil
}
