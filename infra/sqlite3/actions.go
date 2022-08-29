package sqlite3

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
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

// GetSellOrderBooks returns the list of sell order books for a sell order ID
// TODO: Missing tests
func (db Database) GetSellOrderBooks(sellOrderID uuid.UUID) ([]domain.SellOrderBook, error) {
	stm, err := db.Connection.Prepare("SELECT * FROM  sell_order_books WHERE sell_order_id = ?")
	if err != nil {
		return nil, fmt.Errorf("sqlite error: %+v", err)
	}
	defer stm.Close()

	rows, err := stm.Query(sellOrderID.String())
	if err != nil {
		log.Fatal(err)
	}

	var sellOrderBooks []domain.SellOrderBook

	// These paramenters stored in the DB are not useful for the domain side
	var id int
	var createdAt time.Time

	for rows.Next() {
		var e domain.SellOrderBook
		if err := rows.Scan(&id, &e.SellOrderID, &e.BidID,
			&e.Symbol, &e.Quantity, &e.MinimumSellPrice, &e.BidPrice, &createdAt); err != nil {
			return nil, fmt.Errorf("sqlite error: %+v", err)
		}
		sellOrderBooks = append(sellOrderBooks, e)
	}

	return sellOrderBooks, nil
}
