package sqlite3

import (
	"fmt"

	"github.com/jkmrto/trade_executor/domain"
)

// ApplySell stores in the database an opportunity to appy a sell
// TODO: Missing tests
func (db Database) ApplySell(sb domain.SellOrderBook) error {
	stmt, err := db.Connection.Prepare(`
		INSERT INTO sell_order_books(
			sell_order_id,
			bid_id,
			symbol,
			quantity,
			minimum_sell_price,
			bid_price
		) VALUES(?,?,?,?,?,?)`)

	if err != nil {
		return fmt.Errorf("sqlite error: %+v", err)
	}

	defer stmt.Close()

	if _, err = stmt.Exec(sb.SellOrderID, sb.BidID, sb.Symbol, sb.Quantity, sb.MinimumSellPrice, sb.BidPrice); err != nil {
		return err
	}

	return nil
}
