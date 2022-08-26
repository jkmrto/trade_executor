# Trade Executor

## Roadmap

[x] Domain
	- [x] SellOrder (Price, Amount)
[ ] Use cases on the app layer
	- ProcessBidHandler()
[ ] Infra layer:
  - [x] Binance Setup
  - [ ] Database for books order sells
  - [ ] Maybe an HTTP client for setting orders (?)

## Decisions

One of the decison was if processing the bid updates in batches or try to process them one by one. In order to process then OneByOne we would to create one goroutine per each message in order to not block sender channel.

The `wsDepthHandler` receives synchronous call, what means we dont receive any more update if the receiver gets blocked
