mocks:
	moq -out app/zmock_exchange_test.go -pkg app_test app Exchange 

run:
	go run cmd/main.go

test:
	go test  ./... -cover

sell-order-examle1:
	curl --location --request POST 'http://localhost:8080/SellOrder' \
	--header 'Content-Type: application/json' \
	--data-raw '{"price": 300.0, "quantity": 300.0}'
