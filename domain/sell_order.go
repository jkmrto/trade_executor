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
func (so *SellOrder) ApplySellOnBid(bid Bid) *SellBook {
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

	return &SellBook{
		BidID:            bid.ID,
		BidPrice:         bid.Price,
		Symbol:           bid.Symbol,
		MinimumSellPrice: so.Price,
		Quantity:         qty,
	}
}

//SellBook is a register of a sell over a bid
type SellBook struct {
	BidID            string
	BidPrice         float64
	Symbol           string
	MinimumSellPrice float64
	Quantity         float64
}
