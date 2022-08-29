package app

import (
	"github.com/google/uuid"
	"github.com/jkmrto/trade_executor/domain"
)

// Exchange defines the contract for interacting with an exchange
type Exchange interface {
	ApplySell(domain.SellOrderBook) error
	GetSellOrderBooks(uuid.UUID) ([]domain.SellOrderBook, error)
}
