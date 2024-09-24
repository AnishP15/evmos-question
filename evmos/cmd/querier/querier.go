package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// Querier contains fields used for querying the local node and storing the results
type Querier struct {
	rpcClient      *rpc.Client
	ethClient      *ethclient.Client
	walletBalances []WalletBalance
	contracts      map[string]int
	blocks         []*types.Block
}

// WalletBalance stores the Addess and Balance for a wallet
type WalletBalance struct {
	Address string
	Balance *big.Int
}

// NewQuerier instantiates a new Querier connected to an rpc endpoint
func NewQuerier(rpcURL string) (*Querier, error) {
	rpcClient, err := rpc.Dial(rpcURL)
	if err != nil {
		return nil, err
	}

	ethClient, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}
	return &Querier{
		rpcClient:      rpcClient,
		ethClient:      ethClient,
		walletBalances: []WalletBalance{},
		contracts:      make(map[string]int),
		blocks:         []*types.Block{},
	}, nil
}

// storeWalletsAndContracts persists to Querier's mapping of wallet to
// balance based on txs in the stored blocks & contract interaction
func (querier *Querier) storeWalletsAndContracts() error {
	// Create a set-like map to track processed addresses
	walletCache := make(map[string]struct{})

	for _, block := range querier.blocks {
		fmt.Printf("Aggregating statistics for block %v \n", block.Number())
		for _, tx := range block.Transactions() {

			err := querier.storeWallet(tx, walletCache)

			if err != nil {
				return err
			}

			err = querier.storeContracts(tx.Hash())

			if err != nil {
				return err
			}
		}
	}
	return nil
}

// isContractAddress determines if an address is a contract or an EOA
func (querier *Querier) isContractAddress(address common.Address) (bool, error) {
	code, err := querier.ethClient.CodeAt(context.Background(), address, nil) // nil is for latest block
	if err != nil {
		return false, err
	}
	return len(code) > 0, nil
}
