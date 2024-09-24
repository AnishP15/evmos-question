package main

import (
	"sort"

	"github.com/ethereum/go-ethereum/common"
)

func (querier *Querier) storeContracts(txHash common.Hash) error {
	// Fetch internal contract calls via transaction trace
	trace, err := querier.getTransactionTrace(txHash)
	if err != nil {
		return err
	}

	// Check if "calls" key exists in the trace
	calls, ok := trace["calls"].([]interface{})
	if !ok || calls == nil {
		// If there are no internal calls, check direct "to" addresses
		if toAddrStr, exists := trace["to"].(string); exists {
			toAddr := common.HexToAddress(toAddrStr)
			isContract, err := querier.isContractAddress(toAddr)
			if err != nil {
				return err
			}
			if isContract {
				querier.contracts[toAddr.Hex()]++
			}
		}

		return nil
	} else {
		// If "calls" is not empty, process the internal contract calls
		for _, call := range calls {
			callMap, ok := call.(map[string]interface{})
			if !ok {
				continue
			}

			if to, exists := callMap["to"]; exists {
				toAddr := common.HexToAddress(to.(string))
				isContract, err := querier.isContractAddress(toAddr)
				if err != nil {
					return err
				}
				if isContract {
					querier.contracts[toAddr.Hex()]++
				}
			}
		}
	}

	return nil
}

func (querier *Querier) sortContractsByInteractions() []string {
	// Extract all contract addresses from the map
	var contractAddresses []string
	for addr := range querier.contracts {
		contractAddresses = append(contractAddresses, addr)
	}

	// Sort the contract addresses by interaction count in descending order
	sort.Slice(contractAddresses, func(i, j int) bool {
		return querier.contracts[contractAddresses[i]] > querier.contracts[contractAddresses[j]]
	})

	return contractAddresses
}

// getTransactionTrace queries the trace of a transaction for internal contract interactions
func (querier *Querier) getTransactionTrace(txHash common.Hash) (map[string]interface{}, error) {
	var result map[string]interface{}

	// Convert txHash to hex string before using in JSON-RPC call
	txHashHex := txHash.Hex()

	params := []interface{}{txHashHex, map[string]interface{}{"tracer": "callTracer"}}

	// Make the RPC call
	err := querier.rpcClient.Call(&result, "debug_traceTransaction", params...)
	if err != nil {
		return nil, err
	}

	return result, nil
}
