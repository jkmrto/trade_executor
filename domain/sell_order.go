package domain

import (
	"github.com/google/uuid"
)

//SellOrder ...
type SellOrder struct {
	ID                uuid.UUID
	Price             float64
	Symbol            string
	Quantity          float64
	RemainingQuantity float64
}

// NewSellOrder is a constructor
// Missing tests for this constructor
func NewSellOrder(symbol string, price, quantity float64) SellOrder {
	return SellOrder{
		ID:                uuid.New(),
		Symbol:            symbol,
		Price:             price,
		Quantity:          quantity,
		RemainingQuantity: quantity,
	}
}

//CanBidBeApplied ...
func (so SellOrder) canBidBeAccepted(bid Bid) bool {
	return bid.Price >= so.Price
}

// ApplySellOnBid ...
func (so *SellOrder) ApplySellOnBid(bid Bid) *SellOrderBook {
	if !so.canBidBeAccepted(bid) {
		return nil
	}

	var qty float64
	if bid.Quantity > so.RemainingQuantity {
		qty = so.RemainingQuantity
		so.RemainingQuantity = 0
	} else {
		so.RemainingQuantity = so.RemainingQuantity - bid.Quantity
		qty = bid.Quantity
	}

	return &SellOrderBook{
		SellOrderID:      so.ID,
		BidID:            bid.ID,
		BidPrice:         bid.Price,
		Symbol:           bid.Symbol,
		MinimumSellPrice: so.Price,
		Quantity:         qty,
	}
}

//SellOrderBook is a register of a sell over a bid
type SellOrderBook struct {
	SellOrderID      uuid.UUID
	BidID            string
	Symbol           string
	Quantity         float64
	MinimumSellPrice float64
	BidPrice         float64
}
