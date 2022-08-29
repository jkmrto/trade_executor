mocks:
	moq -out app/zmock_exchange_test.go -pkg app_test app Exchange 

run:
	go run cmd/main.go

test:
	go test  ./... -cover

sell-order-example-BNBUSDT:
	curl --location --request POST 'http://localhost:8080/SellOrder' \
	--header 'Content-Type: application/json' \
	--data-raw '{"price": 270.0, "quantity": 300.0, "symbol": "BNBUSDT"}'

sell-order-example-BTCUSDT:
	curl --location --request POST 'http://localhost:8080/SellOrder' \
	--header 'Content-Type: application/json' \
	--data-raw '{"price": 20100, "quantity": 300.0, "symbol": "BTCUSDT"}'

sell-order-example-ETHBTC:
	curl --location --request POST 'http://localhost:8080/SellOrder' \
	--header 'Content-Type: application/json' \
	--data-raw '{"price": 0.071, "quantity": 300.0, "symbol": "ETHBTC"}'



tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

lint:
	golangci-lint run
