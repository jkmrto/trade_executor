mocks:
	moq -out app/zmock_exchange_test.go -pkg app_test app Exchange 

run:
	go run cmd/main.go

test:
	go test -v ./...

