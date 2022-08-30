package app

import (
	"fmt"

	"github.com/jkmrto/trade_executor/domain"
)

// ShowSellOrderSummaryHandler covers the use case of showing a summary for an order that was sold
type ShowSellOrderSummaryHandler struct {
	Exchange Exchange
}

// NewShowSellOrderSummaryHandler  is a constructor
func NewShowSellOrderSummaryHandler(exchange Exchange) ShowSellOrderSummaryHandler {
	return ShowSellOrderSummaryHandler{
		Exchange: exchange,
	}
}

// Handle will do a sell of an order if the bid price is good enough
func (ssosh ShowSellOrderSummaryHandler) Handle(sellOrder domain.SellOrder) error {
	sellOrderBooks, err := ssosh.Exchange.GetSellOrderBooks(sellOrder.ID)
	if err != nil {
		return fmt.Errorf("Error when getting sell order books: %v", err)
	}

	fmt.Printf("\n--------------------------------------------------------")
	fmt.Printf("\norderID\t\t\t\t\tRequiredPrice\tSellPrice\tQuantity\n")
	for _, sob := range sellOrderBooks {
		fmt.Printf("%s\t%+v\t\t%+v\t\t%+v\n", sob.SellOrderID.String(), sob.MinimumSellPrice, sob.BidPrice, sob.Quantity)

	}

	fmt.Printf("--------------------------------------------------------\n")
	return nil
}
