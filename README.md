# ETH LogCache

## Usage

### Start the service using Docker Compose

```
docker-compose up -d
```

or, you can build and run with environment variables:

```
RPC_URL=https://eth-mainnet.reddio.com START_BLOCK=15560257 go run .
```

### Send request to ETH LogCache

```sh
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "method": "eth_getLogs",
    "params": [
      {
        "fromBlock": "0x0",
        "toBlock": "0x120dc53",
        "address": "0xb62bcd40a24985f560b5a9745d478791d8f1945c",
        "topics": [
          [
            "0xcfb473e6c03f9a29ddaf990e736fa3de5188a0bd85d684f5b6e164ebfbfff5d2"
          ]
        ]
      }
    ],
    "id": 62,
    "jsonrpc": "2.0"
  }' \
  http://localhost:3000
```

Response:
```json
{"block_number":[{"id":62,"jsonrpc":"2.0","result":[{"address":"0xb62bcd40a24985f560b5a9745d478791d8f1945c","blockHash":"0xb1c05e3a5f7791b40d9ded2bb67bd2d250f1ccb036dac0f6a046b7ed2d416df0","blockNumber":"0xed6e42","data":"0x0000000000000000000000006b7763b749073e892c83e674c1ec4799d6f339ef","logIndex":"0x151","removed":false,"topics":["0xcfb473e6c03f9a29ddaf990e736fa3de5188a0bd85d684f5b6e164ebfbfff5d2"],"transactionHash":"0x87130dfe52f1eb4ec22261333534ee7ac2c15e5256ffa7a59ae7153119c6cd73","transactionIndex":"0xc9"},{"address":"0xb62bcd40a24985f560b5a9745d478791d8f1945c","blockHash":"0xb1c05e3a5f7791b40d9ded2bb67bd2d250f1ccb036dac0f6a046b7ed2d416df0","blockNumber":"0xed6e42","data":"0x0000000000000000000000006ea99c6fe2c770c2c46ebe03a4855977282e844f","logIndex":"0x153","removed":false,"topics":["0xcfb473e6c03f9a29ddaf990e736fa3de5188a0bd85d684f5b6e164ebfbfff5d2"],"transactionHash":"0xca00caf8ad277bd23597d497e1c77b4308d40e8787d7e0e5204d320f1f3ab31c","transactionIndex":"0xcb"}]}]}
```

## License

MIT