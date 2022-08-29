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

		bidsRouter := app.NewBidsRouter(symbol)
		go bidsRouter.Start()

		sellOrder := domain.NewSellOrder(symbol, 280.0, 100.0)

		sellOrderExecutor := app.NewSellOrderExecutor(&sellOrder, app.ProcessBidHandler{}, app.ShowSellOrderSummaryHandler{}, bidsRouter.SoExecutorFinishedIDCh)
		bidsRouter.NewSellOrderExecutorCh <- &sellOrderExecutor

		for {
			if len(bidsRouter.SoExecutors) == 1 {
				break
			}
		}

		// Check current status of the SellOrderManager
		require.Len(t, bidsRouter.SoExecutors, 1)
		require.Equal(t, bidsRouter.SoExecutors[0].ID, sellOrder.ID)
	})

	t.Run(`Given a bids router with an active sell order manager,
	when a finish sell order message arrived,
	then it removes the sell order manager from the list`, func(t *testing.T) {

		bidsRouter := app.NewBidsRouter(symbol)
		go bidsRouter.Start()

		sellOrder := domain.NewSellOrder(symbol, 280.0, 100.0)

		sellOrderExecutor := app.NewSellOrderExecutor(&sellOrder, app.ProcessBidHandler{}, app.ShowSellOrderSummaryHandler{}, bidsRouter.SoExecutorFinishedIDCh)
		bidsRouter.NewSellOrderExecutorCh <- &sellOrderExecutor

		for {
			if len(bidsRouter.SoExecutors) == 1 {
				break
			}
		}

		bidsRouter.SoExecutorFinishedIDCh <- sellOrderExecutor.ID

		for {
			if len(bidsRouter.SoExecutors) == 0 {
				break
			}
		}

		require.Len(t, bidsRouter.SoExecutors, 0)

	})

	t.Run(`Given a bids router with two sell order manager actived
	when a new batch of bids arrived,
	then both the order managers process the bids`, func(t *testing.T) {

		bidsRouter := app.NewBidsRouter(symbol)
		go bidsRouter.Start()

		sellOrder1 := domain.NewSellOrder(symbol, 280.0, 100.0)
		sellOrder2 := domain.NewSellOrder(symbol, 280.0, 100.0)

		applySellOnExchangeCh := make(chan struct{})
		exchangeMock := &ExchangeMock{
			ApplySellFunc: func(domain.SellOrderBook) error {
				applySellOnExchangeCh <- struct{}{}
				return nil
			},
			GetSellOrderBooksFunc: func(uuid.UUID) ([]domain.SellOrderBook, error) {
				return nil, nil
			},
		}

		pbHandler := app.ProcessBidHandler{Exchange: exchangeMock}
		showSellOrderSummaryHandler := app.NewShowSellOrderSummaryHandler(exchangeMock)

		sellOrderExecutor1 := app.NewSellOrderExecutor(&sellOrder1, pbHandler, showSellOrderSummaryHandler, bidsRouter.SoExecutorFinishedIDCh)
		go sellOrderExecutor1.ProcessBids()
		sellOrderExecutor2 := app.NewSellOrderExecutor(&sellOrder2, pbHandler, showSellOrderSummaryHandler, bidsRouter.SoExecutorFinishedIDCh)
		go sellOrderExecutor2.ProcessBids()

		bidsRouter.NewSellOrderExecutorCh <- &sellOrderExecutor1
		bidsRouter.NewSellOrderExecutorCh <- &sellOrderExecutor2

		for {
			if len(bidsRouter.SoExecutors) == 2 {
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
		require.Equal(t, exchangeMock.ApplySellCalls()[0].SellOrderBook.BidID, bid.ID)
		require.Equal(t, exchangeMock.ApplySellCalls()[1].SellOrderBook.BidID, bid.ID)

		// TODO: Maybe this test should be calling an inteface over
		// ProcessBidHandler, that way we decouple even more the bidsRouter from the exchange
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
