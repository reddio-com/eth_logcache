package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/cockroachdb/pebble"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gofiber/fiber/v2"
)

var (
	RPC_URL       = "https://eth-mainnet.reddio.com"
	ETH_CLIENT, _ = ethclient.Dial(RPC_URL)
	err           error
	START_BLOCK   = uint64(1)
	HTTP_CLIENT   = &http.Client{
		Timeout: 10 * time.Second,
	}
)

func init() {

	// Read from ENV for override
	if os.Getenv("RPC_URL") != "" {
		RPC_URL = os.Getenv("RPC_URL")
	}
	// Read from ENV for override
	if os.Getenv("START_BLOCK") != "" {
		START_BLOCK_STR := os.Getenv("START_BLOCK")
		START_BLOCK_INT, _ := strconv.Atoi(START_BLOCK_STR)
		START_BLOCK = uint64(START_BLOCK_INT)
	}

	DB, err = pebble.Open("data", &pebble.Options{})
	if err != nil {
		log.Fatal(err)
	}

	// Check if LatestBlockNumberInDB exists, if not, create it
	LatestBlockNumberInDBKey, err := LatestBlockNumberInDB()
	if err != nil {
		log.Fatal(err)
	}
	if LatestBlockNumberInDBKey == 0 {
		fmt.Println("LatestBlockNumberInDB not found, creating it")
		SetKeyNumber("LatestBlockNumberInDB", START_BLOCK)
	}
}

func main() {
	// Init runner
	go runner()

	// Init API
	app := fiber.New()

	app.Post("/", payload_handler)
	app.Get("/latest_block", get_latest_block)
	app.Get("/get_block_number", get_block_number)

	app.Listen(":3000")
}
