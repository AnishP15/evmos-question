package main

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

// storeBlocks persists to Querier's blocks field the blocks between two indices
func (querier *Querier) storeBlocks(start int64, end int64) error {
	var blocks []*types.Block
	for i := start; i <= end; i++ {
		block, err := querier.getBlock(i)

		if err != nil {
			return err
		}

		blocks = append(blocks, block)
	}
	querier.blocks = blocks
	return nil
}

// getBlock is a helper function that queries for a block at a specific height
func (querier *Querier) getBlock(blockNumber int64) (*types.Block, error) {
	blockNum := big.NewInt(blockNumber)
	block, err := querier.ethClient.BlockByNumber(context.Background(), blockNum)

	if err != nil {
		return nil, err
	}

	return block, nil
}
