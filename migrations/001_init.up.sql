CREATE TABLE sell_order_books (
  id INTEGER PRIMARY KEY,
  sell_order_id TEXT NOT NULL,
  bid_id TEXT NOT NULL,
  symbol TEXT NOT NULL,
  quantity REAL NOT NULL,
  minimum_sell_price REAL NOT NULL,
  bid_price REAL NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

