package main

import (
	"context"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func (querier *Querier) storeWallet(tx *types.Transaction, walletCache map[string]struct{}) error {
	chainID, err := querier.ethClient.NetworkID(context.Background())

	if err != nil {
		return err
	}

	// Extract 'from' address
	msg, err := tx.AsMessage(types.NewEIP155Signer(chainID), tx.GasPrice())
	if err != nil {
		return err
	}
	from := msg.From().Hex()

	// Check if 'from' address has been processed, and if not, add to walletBalances
	if _, exists := walletCache[from]; !exists {

		isContract, _ := querier.isContractAddress(msg.From())

		if !isContract {

			walletCache[from] = struct{}{}
			querier.walletBalances = append(querier.walletBalances, WalletBalance{
				Address: from,
				Balance: common.Big0,
			})
		}
	}

	// Process the 'to' address, if it exists
	to := tx.To()
	if to != nil {
		toAddress := to.Hex()

		// Check if 'to' address has been processed, and if not, add to walletBalances
		if _, exists := walletCache[toAddress]; !exists {

			isContract, _ := querier.isContractAddress(*to)

			if !isContract {
				walletCache[toAddress] = struct{}{}
				querier.walletBalances = append(querier.walletBalances, WalletBalance{
					Address: toAddress,
					Balance: common.Big0,
				})
			}
		}
	}
	return nil
}

// populateWalletBalances queries the balance for each wallet address stored
func (querier *Querier) populateWalletBalances() error {

	for i := range querier.walletBalances {
		bal, err := querier.getWalletBalance(querier.walletBalances[i].Address)

		if err != nil {
			return nil
		}
		querier.walletBalances[i].Balance = bal
	}

	return nil
}

// getWalletBalance returns the balances associated with a wallet
func (querier *Querier) getWalletBalance(address string) (*big.Int, error) {
	account := common.HexToAddress(address)
	balance, err := querier.ethClient.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

// sortWalletsByBalance sorts wallet addresses by their balances in descending order
func (querier *Querier) sortWalletsByBalance() {
	sort.Slice(querier.walletBalances, func(i, j int) bool {
		// Compare balances to sort in descending order
		return querier.walletBalances[i].Balance.Cmp(querier.walletBalances[j].Balance) > 0
	})
}
