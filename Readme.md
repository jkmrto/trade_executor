# Trade Executor

## Roadmap

[x] Domain
	- [x] SellOrder (Price, Amount)
[ ] Use cases on the app layer
	- ProcessBidHandler()
[ ] Infra layer:
  - [x] Binance Setup 
	-
	- Do we need a `BidUpdatesRouter` (?). This manager may be in charge in sending the Bid updates to each 
  - [ ] Interface for exchange actions
		- [ ] It can be a `sqlite3` client for storing the result.
		- [ ] It can be just an in memory solution.
  - [ ] Maybe an HTTP client for setting orders (?). 


## Design decisions

- What about having a goroutine per each sell order? That way we could process in parallel the bid updates, so we can take a decision as soon as possible. As first approach I will assume only one sell order, but I may have to go in this direction later. 



## Decisions

One of the decison was if processing the bid updates in batches or try to process them one by one. In order to process then OneByOne we would to create one goroutine per each message in order to not block sender channel.

The `wsDepthHandler` receives synchronous call, what means we dont receive any more update if the receiver gets blocked
