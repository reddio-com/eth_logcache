package main

import (
	"fmt"
)

func runner() {
	for {
		// Get LatestBlockNumberOnChain
		LatestBlockNumberOnChain, err := LatestBlockOnChain()
		if err != nil {
			fmt.Printf("Failed to get latest block number on chain: %v\n", err)
			return
		}
		fmt.Printf("Latest block number on chain: %d\n", LatestBlockNumberOnChain)
		// Get Key with LatestBlockNumberInDB
		LatestBlockNumberInDBKey, err := LatestBlockNumberInDB()
		if err != nil {
			fmt.Printf("Failed to get latest block number in DB: %v\n", err)
			return
		}

		// blockNumber := 19088166
		// value := make([]byte, 8)
		// binary.LittleEndian.PutUintgo64(value, uint64(blockNumber))
		// fmt.Println("Value: ", value)

		// Get block number from DB
		fmt.Printf("Latest block number in DB: %d\n", LatestBlockNumberInDBKey)

		if LatestBlockNumberInDBKey < LatestBlockNumberOnChain {
			for i := LatestBlockNumberInDBKey; i < LatestBlockNumberOnChain; i++ {
				fmt.Println("Dealing with block number: ", i)
				// Get topic list from block number
				topicList, err := getTopicListByBlockNumber(i)
				if err != nil {
					fmt.Printf("Failed to get topic list: %v\n", err)
					// Retry
					return
				}
				fmt.Printf("Total topics: %d\n", len(topicList))
				if len(topicList) > 0 {
					for _, topic := range topicList {
						SetKeyNumberArray(topic, i)
					}
					// Flush DB
					FlushDB()
				}
				// Update LatestBlockNumberInDB
				SetKeyNumber("LatestBlockNumberInDB", i)
			}
		}
		defer DB.Close()
	}
}
