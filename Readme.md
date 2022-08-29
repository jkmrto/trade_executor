# Trade Executor

You can find a makefile with some useful command for the project. 

- Running it:

```
make run
```

- Creating a sell order:

```
sell-order-example-BTCUSDT
sell-order-example-BNBUSDT
sell-order-example-ETHBTC
```

- Run the test

```
make test
```

## How did you approach the task?

I tried to follow a hexagonal architecture approach, trying to define the domain and the use cases. In a hexagonal architecture, we divide the application into three layers. 

- **Domain**. Only the definition of the entities and business logic without side effects.
- **App**. Use cases for the application. 
- **Infra**. Implementation for communicating for external services/clients.  

I wanted to only have one Binance consumer per symbol, so I built a pipeline for processing the bids per symbol.

- **BinanceListener**. It listens to the updates, sending them to his **BidsRouter**. We will have one listener per symbol. The number of **BinanceListener** is determined by the number of supported symbols.

- **BidsRouter**: It listens to the bid updates for the Binance consumer. It routes these updates to all **SellOrderExecutor**. It is executed in a `goroutine`.

- **SellOrderExecutor**: It receives the bid updates from a **BidsRouter**, checking if the SellOrder can be applied. It is launched in a `goroutine` that will finish once the order is executed. We will have as many SellOrderExecutor as orders are in the process of being sold in the system.


Out of this pipeline,  we have the **SellOrderManager**. It receives the new sell order requests from the HTTP server. On each request,  It launches a new  **SellOrderExecutor**, connecting it to his respective BidsRouter.


All these components are not easy to classify in a hexagonal architecture. That makes me think that trying to follow this architecture maybe was not the best approach.

### Another possible approach

 Just launching a Binance consumer per each SellOrder. I see three bad points here:
- we may be overloading the Binance publisher.
- We could be banned if we create too many connections to the Binance WS.
- We will be not making the best use of the resources of the system.


## Where did you have difficulties?

Trying to build the pipeline for processing the bids but still trying to keep a hexagonal architecture.

## What part did you like the most?

The implementation of the **BidsRouter**. It is an interesting piece since it gets notifications about new **SellOrderExecutors** from the **SellOrderManager**, but it also receives the updates to remove **SellOrderExecutor** when the order is sold.  


## How much time you spent on it in total?

Probably from 9 to 10 hours.

## Ideal architecture -> Microservices 

This implementation is just a basic approach to the problem. In the long term, we could think about splitting the application into several microservices, moving the logic of routing the bids out of the application to a messages broker like RabbitMQ or Kafka. We could think about these services (I will use RabbitMQ terminology since is the tool I am more used to working with)

- **Binance Consumer** -> This service connects with the Binance platform for listening to the bid updates. It publishes the bids into a RabbitMQ exchange.  

- **Sell Order Receiver** -> This service is just a HTTP interface for receiving updates from the clients. It publishes these new sell orders into RabbitMQ

- **Trade Executor** -> This service consumes the Binance updates from RabbitMQ. We will have a queue for each sell Order. This way we move out all the logic for routing the bids out of the application, so we can focus on the use cases.   

We could even think of forms of scaling the application. Maybe, a good approach will be to scale this **Trade Executor** by symbol. So we could have an instance per each symbol


## Main points of improvements

- [ ] The application is not resilient to restarts. Now the application is stateful regarding the orders that are being sold, so if we do a restart we will lose which sell orders are running at the moment. We will need some kind of recovery mechanism to restore the status of the service after a restart. This is quite important for a cloud-native environment. 
- [ ] The shutdown process is not optimal for production. We are not closing the running goroutines, the DB connection and the Binance consumers properly.
- [ ] Tests missing on HTTP server, sqlite3 layer and the Binance listener.
- [ ] Share the context along all the components in the application. We should use the context for cancelling active actions (eg. operations on the DB) and the `ctx.Done` channel to close the running goroutines. 
- [ ] Add github actions for pipelines checks
- [ ] Add some kind of logger (?)



## Roadmap

[x] Domain
	- SellOrder (Price, Amount)
	- SellOrderBook
	- Bid
[x] Use cases on the app layer
	- ProcessBidHandler()
	- ShowSellOrderSummary()
[x] Schedulers/Organizers
  - [x] SellOrderExecutor
  - [x] SellOrderManager
  - [x] BidsRouter
[x] Infra layer:
  - [x] Binance Setup
    - It is configurable by Symbol -> We will have one Binance consumer per symbol.
  - [x] Interface for exchange actions -> sqlite3 implementaion.
    - It is a registry of the sell order operations 
  - [x] HTTP endpoint for setting orders. 
