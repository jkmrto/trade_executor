package app_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/jkmrto/trade_executor/app"
	"github.com/jkmrto/trade_executor/domain"
	"github.com/stretchr/testify/require"
)

func TestTest(t *testing.T) {
	t.Run(`When it is called with a bid that can be accepted,
	then it accepts the bid and send the order to the exchange`,
		func(t *testing.T) {

			exchangeMock := &ExchangeMock{
				ApplySellFunc: func(*domain.SellBook) error {
					fmt.Printf("this is being called")
					return nil
				},
			}

			pbHandler := app.ProcessBidHandler{Exchange: exchangeMock}

			bid := domain.NewBid(uuid.New().String(), "BTCUSDT", 1.0, 1.0)
			so := domain.NewSellOrder("BTCUSDT", 1.0, 1.0)

			pbHandler.Handle(&so, bid)

			require.Len(t, exchangeMock.ApplySellCalls(), 1)
			sellBook := exchangeMock.ApplySellCalls()[0].SellBook
			require.Equal(t, bid.Quantity, sellBook.Quantity)

			require.Equal(t, so.RemainingQuantity, 0.0)

		})

	t.Run(`When it is called with a bid that can not be accepted,
	then id doesn't call the exchange`,
		func(t *testing.T) {

			// We could even remove this mock, since it is not going to be called
			// But let's keep it just for being explicits on this test
			exchangeMock := &ExchangeMock{
				ApplySellFunc: func(*domain.SellBook) error {
					fmt.Printf("this is being called")
					return nil
				},
			}

			pbHandler := app.ProcessBidHandler{Exchange: exchangeMock}

			sellOrderQty := 10.0
			bid := domain.NewBid(uuid.New().String(), "BTCUSDT", 1.0, 1.0)
			so := domain.NewSellOrder("BTCUSDT", 10.0, sellOrderQty)

			pbHandler.Handle(&so, bid)
			require.Empty(t, exchangeMock.ApplySellCalls(), 1)
			require.Equal(t, so.RemainingQuantity, sellOrderQty)

		})

}
