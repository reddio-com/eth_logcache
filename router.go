package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Payload struct {
	Jsonrpc string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	ID      int      `json:"id"`
	Params  []Params `json:"params"`
}

type Params struct {
	Topics    []TopicList `json:"topics"`
	FromBlock string      `json:"fromBlock"`
	ToBlock   string      `json:"toBlock"`
	Address   string      `json:"address"`
}

type TopicList []string

// Given [1,2,3,5,6,7], return [[1,3],[5,7]]
func getListOfStartAndToBlockNumber(blockNumberList []uint64) [][]uint64 {
	startBlockNumber := blockNumberList[0]
	toBlockNumber := blockNumberList[0]
	blockNumberRangeList := make([][]uint64, 0)
	for i := 1; i < len(blockNumberList); i++ {
		if blockNumberList[i] == toBlockNumber+1 {
			toBlockNumber = blockNumberList[i]
		} else {
			blockNumberRangeList = append(blockNumberRangeList, []uint64{startBlockNumber, toBlockNumber})
			startBlockNumber = blockNumberList[i]
			toBlockNumber = blockNumberList[i]
		}
	}
	blockNumberRangeList = append(blockNumberRangeList, []uint64{startBlockNumber, toBlockNumber})
	return blockNumberRangeList
}

// Payload handler is for user facing handler
// It will check if the method is eth_getLogs and use Pebble to get block_numbers
// Other requests will be forwarded to RPC_URL

// Example payload
//
//	{
//		"method": "eth_getLogs",
//		"params": [
//		  {
//			"fromBlock": "0x0",
//			"toBlock": "0x1234325",
//			"address": "0xb62bcd40a24985f560b5a9745d478791d8f1945c",
//			"topics": [
//			  [
//				"0xcfb473e6c03f9a29ddaf990e736fa3de5188a0bd85d684f5b6e164ebfbfff5d2"
//			  ]
//			]
//		  }
//		],
//		"id": 62,
//		"jsonrpc": "2.0"
//	  }
func payload_handler(c *fiber.Ctx) error {
	// Check method is eth_getLogs
	payload := new(Payload)
	if err := c.BodyParser(payload); err != nil {
		return fmt.Errorf("Failed to parse body: %v\n", err)
	}
	if payload.Method == "eth_getLogs" {
		// Check if topics is empty
		if len(payload.Params[0].Topics) > 0 {

			// TODO: Here only processeed the first topic, need to process all topics
			topicList0 := payload.Params[0].Topics[0]
			combinedBlockNumberList := make([]uint64, 0)
			for _, topic := range topicList0 {
				blockNumberArray, err := GetKeyNumberArray(topic)
				if err != nil {
					return fmt.Errorf("Failed to get block number: %v\n", err)
				}
				combinedBlockNumberList = append(combinedBlockNumberList, blockNumberArray...)
			}

			// Check if empty
			if len(combinedBlockNumberList) == 0 {
				return c.JSON(fiber.Map{
					"block_number": combinedBlockNumberList,
				})
			}

			// Dedup and sort
			combinedBlockNumberList = DeduplicateUint64Array(combinedBlockNumberList)
			sort.Slice(combinedBlockNumberList, func(i, j int) bool {
				return combinedBlockNumberList[i] < combinedBlockNumberList[j]
			})

			// BlockNumber Ranges
			blockNumberRanges := getListOfStartAndToBlockNumber(combinedBlockNumberList)

			// Now we have a mixed block number ranges, we do requests to RPC_URL
			// Make sure to check if the block number is in the FromBlock and ToBlock
			FromBlockDec, err := strconv.ParseUint(payload.Params[0].FromBlock, 0, 64)
			if err != nil {
				return fmt.Errorf("Failed to parse fromBlock: %v\n", err)
			}
			ToBlockDec, err := strconv.ParseUint(payload.Params[0].ToBlock, 0, 64)
			if err != nil {
				return fmt.Errorf("Failed to parse toBlock: %v\n", err)
			}

			combinedResponse := make([]map[string]interface{}, 0)

			for _, eachBlockRange := range blockNumberRanges {

				if eachBlockRange[0] < FromBlockDec {
					eachBlockRange[0] = FromBlockDec
				}
				if eachBlockRange[1] > ToBlockDec {
					eachBlockRange[1] = ToBlockDec
				}
				if eachBlockRange[0] > eachBlockRange[1] {
					// swap
					eachBlockRange[0], eachBlockRange[1] = eachBlockRange[1], eachBlockRange[0]
				}

				// Do request to RPC_URL
				newPayload := &Payload{
					Jsonrpc: payload.Jsonrpc,
					Method:  payload.Method,
					ID:      payload.ID,
					Params: []Params{
						{
							Topics:    payload.Params[0].Topics,
							FromBlock: fmt.Sprintf("0x%x", eachBlockRange[0]),
							ToBlock:   fmt.Sprintf("0x%x", eachBlockRange[1]),
							Address:   payload.Params[0].Address,
						},
					},
				}
				fmt.Printf("New payload: %+v\n", newPayload)
				payloadByteSlice, err := json.Marshal(newPayload)
				if err != nil {
					return fmt.Errorf("Failed to marshal payload: %v\n", err)
				}

				// Send request to RPC_URL
				response, err := HTTP_CLIENT.Post(RPC_URL, "application/json", bytes.NewBuffer(payloadByteSlice))
				if err != nil {
					return fmt.Errorf("Failed to send request to RPC_URL: %v\n", err)
				}
				defer response.Body.Close()

				// Read response
				var result map[string]interface{}
				json.NewDecoder(response.Body).Decode(&result)
				combinedResponse = append(combinedResponse, result)
				fmt.Printf("Result: %+v\n", result)
			}

			return c.JSON(fiber.Map{
				"block_number": combinedResponse,
			})
		}
	}
	return nil
}

func get_latest_block(c *fiber.Ctx) error {
	LatestBlockNumberInDBKey, err := LatestBlockNumberInDB()
	if err != nil {
		return fmt.Errorf("Failed to get latest block number in DB: %v\n", err)
	}
	return c.JSON(fiber.Map{
		"latest_block": LatestBlockNumberInDBKey,
	})
}

func get_block_number(c *fiber.Ctx) error {
	topic := c.Query("topic")
	if topic == "" {
		return fmt.Errorf("Topic is empty")
	}
	blockNumberArray, err := GetKeyNumberArray(topic)
	if err != nil {
		return fmt.Errorf("Failed to get block number: %v\n", err)
	}
	if len(blockNumberArray) == 0 {
		return c.JSON(fiber.Map{
			"block_number": blockNumberArray,
		})
	}
	// Dedup and sort
	blockNumberArray = DeduplicateUint64Array(blockNumberArray)
	sort.Slice(blockNumberArray, func(i, j int) bool {
		return blockNumberArray[i] < blockNumberArray[j]
	})

	// BlockNumber Ranges
	blockNumberRanges := getListOfStartAndToBlockNumber(blockNumberArray)

	return c.JSON(fiber.Map{
		"block_number": blockNumberRanges,
	})
}
