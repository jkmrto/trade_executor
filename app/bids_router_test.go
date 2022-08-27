package app_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jkmrto/trade_executor/app"
	"github.com/jkmrto/trade_executor/domain"
	"github.com/stretchr/testify/require"
)

const symbol = "BNBUSDT"

func TestBidsRouter(t *testing.T) {

	t.Run(`Given a bids router wihtout any active sell order manager,
	when a new sell order manager message arrived,
	then it adds this sell order manager to the list`, func(t *testing.T) {

		bidsRouter := app.NewBidsRouter()
		go bidsRouter.Start()

		sellOrder := domain.NewSellOrder(symbol, 280.0, 100.0)

		sellOrderManager := app.NewSellOrderManager(&sellOrder, app.ProcessBidHandler{}, bidsRouter.SoManagerFinishedIDCh)
		bidsRouter.NewSellOrderManagerCh <- &sellOrderManager

		for {
			if len(bidsRouter.SoManagers) == 1 {
				break
			}
		}

		// Check current status of the SellOrderManager
		require.Len(t, bidsRouter.SoManagers, 1)
		require.Equal(t, bidsRouter.SoManagers[0].ID, sellOrderManager.ID)
	})

	t.Run(`Given a bids router with an active sell order manager,
	when a finish sell order message arrived,
	then it removes the sell order manager from the list`, func(t *testing.T) {

		bidsRouter := app.NewBidsRouter()
		go bidsRouter.Start()

		sellOrder := domain.NewSellOrder(symbol, 280.0, 100.0)

		sellOrderManager := app.NewSellOrderManager(&sellOrder, app.ProcessBidHandler{}, bidsRouter.SoManagerFinishedIDCh)
		bidsRouter.NewSellOrderManagerCh <- &sellOrderManager

		for {
			if len(bidsRouter.SoManagers) == 1 {
				break
			}
		}

		bidsRouter.SoManagerFinishedIDCh <- sellOrderManager.ID

		for {
			if len(bidsRouter.SoManagers) == 0 {
				break
			}
		}

		require.Len(t, bidsRouter.SoManagers, 0)

	})

	t.Run(`random`, func(t *testing.T) {

		bidsRouter := app.NewBidsRouter()
		go bidsRouter.Start()

		sellOrder1 := domain.NewSellOrder(symbol, 280.0, 100.0)
		sellOrder2 := domain.NewSellOrder(symbol, 280.0, 100.0)

		applySellOnExchangeCh := make(chan struct{})
		exchangeMock := &ExchangeMock{
			ApplySellFunc: func(domain.SellBook) error {
				applySellOnExchangeCh <- struct{}{}
				return nil
			},
		}

		pbHandler := app.ProcessBidHandler{Exchange: exchangeMock}

		sellOrderManager1 := app.NewSellOrderManager(&sellOrder1, pbHandler, bidsRouter.SoManagerFinishedIDCh)
		go sellOrderManager1.ProcessBids()
		sellOrderManager2 := app.NewSellOrderManager(&sellOrder2, pbHandler, bidsRouter.SoManagerFinishedIDCh)
		go sellOrderManager2.ProcessBids()

		bidsRouter.NewSellOrderManagerCh <- &sellOrderManager1
		bidsRouter.NewSellOrderManagerCh <- &sellOrderManager2

		for {
			if len(bidsRouter.SoManagers) == 2 {
				break
			}
		}

		bid := newBid(280.0, 100.0)
		bidsRouter.BidsCh <- []domain.Bid{bid}

		// Sihce we have two sell order manager
		// the ApplySell on the Exchange is called twice
		<-applySellOnExchangeCh
		<-applySellOnExchangeCh

		require.Len(t, exchangeMock.ApplySellCalls(), 2)
		require.Equal(t, exchangeMock.ApplySellCalls()[0].SellBook.BidID, bid.ID)
		require.Equal(t, exchangeMock.ApplySellCalls()[1].SellBook.BidID, bid.ID)

		// TODO. Maybe this test should be calling an inteface over
		// ProcessBidHandler, that way we decouple even more the
		// bidsRouter from the exchange
	})

}

func newBid(price, quantity float64) domain.Bid {
	return domain.Bid{
		ID:       uuid.New().String(),
		Symbol:   symbol,
		Price:    price,
		Quantity: quantity,
	}

}
