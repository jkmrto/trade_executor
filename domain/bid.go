package domain

import (
	"errors"
	"strconv"
)

// Bid ...
type Bid struct {
	ID       string
	Symbol   string
	Price    float64
	Quantity float64
}

// NewBid ...
func NewBid(bidID, symbol string, price, quantity float64) Bid {
	return Bid{
		ID:       bidID,
		Symbol:   symbol,
		Price:    price,
		Quantity: quantity,
	}
}

// NewBidFromRaw ...
func NewBidFromRaw(bidID, symbol, rawPrice, rawQuantity string) (Bid, error) {
	price, err := strconv.ParseFloat(rawPrice, 64)
	if err != nil {
		return Bid{}, errors.New("Invalid price")
	}

	quantity, err := strconv.ParseFloat(rawQuantity, 64)
	if err != nil {
		return Bid{}, errors.New("Invalid quantity")
	}

	return NewBid(bidID, symbol, price, quantity), nil

}
