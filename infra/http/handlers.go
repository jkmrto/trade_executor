package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/jkmrto/trade_executor/app"
	"github.com/jkmrto/trade_executor/domain"
)

const symbol = "BNBUSDT"

// SellOrder  is a DTO
type SellOrder struct {
	Quantity *float64 `json:"quantity"`
	Price    *float64 `json:"price"`
	Symbol   string   `json:"symbol"`
}

// CreateSellOrder launces a new SellOrderManager if the request is correct
// TODO: Missing tests for this module
func CreateSellOrder(somOrganizer app.SellOrderManagerOrganizer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var so SellOrder
		if err := json.NewDecoder(r.Body).Decode(&so); err != nil {
			fmt.Printf("%+v", err)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("Unable to decode your request"))
			return
		}

		if so.Price == nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("price is required in the body request"))
			return
		}

		if so.Quantity == nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("quantity is required in the body request"))
			return
		}

		if so.Symbol == "" {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("symbol is required in the body request"))
			return
		}

		domainSellOrder := domain.SellOrder{
			ID:                uuid.New(),
			Symbol:            so.Symbol,
			Price:             *so.Price,
			Quantity:          *so.Quantity,
			RemainingQuantity: *so.Quantity,
		}

		// TODO: I dont really like  being launching from here the SellOrderManager
		// It would be better to do this in a app.CreateSellOrderHandler
		err := somOrganizer.LaunchNewSellOrderManager(domainSellOrder)

		var unsupportedSymbolError app.UnsupportedSymbolError
		if errors.As(err, &unsupportedSymbolError) {
			w.WriteHeader(http.StatusNotImplemented)
			_, _ = w.Write([]byte("The symbol specified is not supported"))
			return
		}

		msg := fmt.Sprintf("The Sell Order was created with ID: \"%+v\"", domainSellOrder.ID)
		_, _ = w.Write([]byte(msg))
		w.WriteHeader(http.StatusOK)
	}
}
