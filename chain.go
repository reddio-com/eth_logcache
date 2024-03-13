package main

import (
	"context"
	"fmt"
	"math/big"
)

func LatestBlockOnChain() (uint64, error) {
	header, err := ETH_CLIENT.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return 0, fmt.Errorf("Failed to get latest block header: %v", err)
	}

	return header.Number.Uint64(), nil
}

func getTopicListByBlockNumber(blockNumber uint64) ([]string, error) {
	blockNumberBig := new(big.Int)
	blockNumberBig.SetUint64(blockNumber)
	block, err := ETH_CLIENT.BlockByNumber(context.Background(), blockNumberBig)
	for err != nil {
		fmt.Printf("Failed to get block: %v\n", err)
		block, err = ETH_CLIENT.BlockByNumber(context.Background(), blockNumberBig)
	}

	totalTopics := make([]string, 0)

	for _, tx := range block.Transactions() {
		// fmt.Printf("Dealing with transaction: %s\n", tx.Hash().Hex())

		receipt, err := ETH_CLIENT.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			return nil, fmt.Errorf("Failed to get transaction receipt: %v", err)
		}

		for _, log := range receipt.Logs {
			for _, topic := range log.Topics {
				totalTopics = append(totalTopics, topic.Hex())
			}
		}
	}

	return totalTopics, nil
}
