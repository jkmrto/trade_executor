package app

import "github.com/jkmrto/trade_executor/domain"

// Exchange define the contract for selling trades
type Exchange interface {
	ApplySell(*domain.SellBook) error
}
