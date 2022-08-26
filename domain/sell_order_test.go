package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestSellOrder_ApplySellOnBid(t *testing.T) {
	t.Run(`When it is called with a bid price lower than the sell order,
		then it returns nil (since the sell has not happened)`,
		func(t *testing.T) {
			bid := NewBid(uuid.New().String(), "BTCUSDT", 10.0, 10.0)
			so := NewSellOrder("BTCUSDT", 40.0, 40.0)

			require.Nil(t, so.ApplySellOnBid(bid))
		})

	t.Run(`When it is called with a bid price higher than the sell order price,
		and the bid quantity is bigger than the quantity to sell ,
		then the sold quantity is the previous remaining and the remaining now is 0`,
		func(t *testing.T) {
			bid := NewBid(uuid.New().String(), "BTCUSDT", 10.0, 10.0)
			previousRemainingQty := 5.0
			so := NewSellOrder("BTCUSDT", 8.0, previousRemainingQty)

			sellBook := so.ApplySellOnBid(bid)
			require.Equal(t, bid.ID, sellBook.BidID)
			require.Equal(t, bid.Price, sellBook.BidPrice)
			require.Equal(t, so.Symbol, sellBook.Symbol)
			require.Equal(t, previousRemainingQty, sellBook.Quantity)
			require.Equal(t, sellBook.MinimumSellPrice, so.Price)

			require.Equal(t, so.RemainingQuantity, 0.0)
		})

	t.Run(`When it is called with a bid price higher than the sell order price,
		and the quantity to sell is bigger that the bid quantity offered,
		then the sold quantity is the bid quantity and the remaining now is the priveous remaining minus the bid quantity`,
		func(t *testing.T) {
			bid := NewBid(uuid.New().String(), "BTCUSDT", 10.0, 5.0)
			previousRemainingQty := 10.0
			so := NewSellOrder("BTCUSDT", 8.0, previousRemainingQty)

			sellBook := so.ApplySellOnBid(bid)

			require.Equal(t, bid.Quantity, sellBook.Quantity)
			require.Equal(t, so.RemainingQuantity, previousRemainingQty-bid.Quantity)
		})

}
