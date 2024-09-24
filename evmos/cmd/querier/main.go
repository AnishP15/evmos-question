package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

const RPC_URL = "http://localhost:8545"

func main() {

	// Instantiate a New Querier object connected to the local node
	querier, err := NewQuerier(RPC_URL)

	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	startBlockIdx := int64(100)
	endBlockIdx := int64(200)

	// Query for blocks between start idx and end idx inclusive and store them on the Querier
	err = querier.storeBlocks(startBlockIdx, endBlockIdx)

	if err != nil {
		log.Fatalf("Failed to get blocks: %v", err)
	}

	// Find wallets that interacted with the network in these blocks
	err = querier.storeWalletsAndContracts()
	if err != nil {
		log.Fatalf("Failed to query wallets and contracts: %v", err)
	}

	// Query the balance associated with each wallet stored
	err = querier.populateWalletBalances()

	if err != nil {
		log.Fatalf("Failed to populate wallets' balances: %v", err)
	}

	// Sort wallets by descending order
	querier.sortWalletsByBalance()

	// Save querier's list of (wallet address, balance) to a csv file
	err = querier.saveWalletBalancesToCSV("wallet_balances.csv")
	if err != nil {
		log.Fatalf("Failed to save wallet balances to CSV: %v", err)
	}

	// Sort contracts in descending order based on how many times they were accessed
	sortedContracts := querier.sortContractsByInteractions()

	// Save querier's mapping of contract to interactions to a csv file
	err = querier.saveContractsToCSV("contract_interactions.csv", sortedContracts)
	if err != nil {
		log.Fatalf("Failed to save contracts to CSV: %v", err)
	}

}

// Write Querier's wallet balances mapping to a CSV file
func (querier *Querier) saveWalletBalancesToCSV(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header for CSV
	writer.Write([]string{"Wallet Address", "Balance"})

	// Write each wallet's address and balance as a row
	for _, wallet := range querier.walletBalances {
		// Convert balance to string
		balance := wallet.Balance.String()
		err := writer.Write([]string{wallet.Address, balance})
		if err != nil {
			return err
		}
	}

	return nil
}

// Write Querier's mapping of balance to frequency of interactions to a CSV file
func (querier *Querier) saveContractsToCSV(filePath string, sortedContracts []string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	writer.Write([]string{"Contract Address", "Interactions"})

	// Write each contract's address and interactions count
	for _, contract := range sortedContracts {
		interactionCount := strconv.Itoa(querier.contracts[contract])
		err := writer.Write([]string{contract, interactionCount})
		if err != nil {
			return err
		}
	}

	return nil
}
