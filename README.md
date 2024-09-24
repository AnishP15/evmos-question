# Evmos Challenge

## How to Run Code

First, run an evmos local node

```go 
cd evmos
evmosd start
```

Then, find the ```querier``` directory

```go
cd cmd/querier
make run
```

This will start the Querier binary and store the results in csv files.


## Technical Decisions

I decided to run a query process independently to allow for greater customizability. In particular, I added a ```Querier``` struct to store 
the on-chain statistics we query for in memory.

```
type Querier struct {
	rpcClient      *rpc.Client
	ethClient      *ethclient.Client
	walletBalances []WalletBalance
	contracts      map[string]int
	blocks         []*types.Block
}

type WalletBalance struct {
	Address string
	Balance *big.Int
}
```

I am storing wallet balances in a slice rather than a map because it is more efficient to sort a slice rather than an unordered map. Also, we only access a particular balance once so it we cannot realize the gains of constant time lookup which a map would provide us. However, I used a map to represent the contracts interacted with since we will need to access a particular contract multiple times to increment the counter if it is interacted with multiple times. Using a slice would be expensive here because we would need to traverse through it multiple times potentially.

In order to actually make the API calls needed for data aggregation, I utilize ```ethClient``` to access RPC endpoints.

## Next Steps
1) Add a CLI so the user can input their query range
2) Add a db to avoid repeat queries and improve efficiency
3) Add a custom tracer to avoid the bottleneck of ```callTracer```
4) Implement visualizations for better user experience