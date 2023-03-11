Bybit WebSocket Go Library
![GitHub Workflow Status](https://github.com/waxdred/bybit_websocket_go/actions/workflows/go.yml/badge.svg)
![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)
This is a Go library for interacting with the Bybit WebSocket API.

## Table of Contents
- [Installation](#Installation)
- [Usage](#Usage)

# Installation
To use this library in your Go project, you can install it using go get:

```go
go get github.com/waxdred/bybit_websocket_go
```
# Usage
Here is an example main function that demonstrates how to use the bybit_websocket_go library to connect to the Bybit WebSocket API:

```go
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	Bwss "github.com/waxdred/bybit_websocket_go"
)

var (
	apiKey    = "your_api_key"
	apiSecret = "your_api_secret"
)

func main() {
	// Create a new Bybit WebSocket client
	client := Bwss.NewClient()

	// Add a private connection for authenticated requests
	privateConn, err := client.AddConnPrivate(Bwss.Testnet, apiKey, apiSecret)
	if err != nil {
		log.Fatal(err)
	}

	// Subscribe to wallet and position updates
	privateConn.AddPrivateSubs([]string{"wallet", "position"}, handlePrivateData)

	// Add a public connection for market data
	publicConn, err := client.AddConnPublic(Bwss.Perpetual, Bwss.Mainnet)
	if err != nil {
		log.Fatal(err)
	}

	// Subscribe to the BTCUSDT order book
	publicConn.AddPublicSubs([]string{"orderBookL2_25.BTCUSD"}, handlePublicData)

	// Start listening for updates
	go client.Listen()

	// Wait for 10 seconds
	time.Sleep(time.Second * 10)

	// Close the public connection
	publicConn.Close()

	// Wait for a SIGINT signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	// Close the private connection and stop listening
	privateConn.Close()
	client.Stop()
}

func handlePrivateData(conn *Bwss.Connection, msg *Bwss.PrivateData) {
	switch msg.Topic {
	case "wallet":
		wallet := msg.Data.(Bwss.Wallet)
		log.Println(wallet.PrettyFormat())
	case "position":
		position := msg.Data.(Bwss.Position)
		log.Println(position.PrettyFormat())
	}
}

func handlePublicData(conn *Bwss.Connection, msg *Bwss.PublicData) {
	switch msg.Topic {
	case "orderBookL2_25.BTCUSD":
		orderBook := msg.Data.(Bwss.OrderBook)
		fmt.Println(orderBook.PrettyFormat())
	}
}
```
In this example, we first import the bybit_websocket_go library:

```go
import (
    Bwss "github.com/waxdred/bybit_websocket_go"
)
```
Next, we define the apiKey and apiSecret variables with your Bybit API key and secret:

```go
var (
    apiKey    = "your_api_key"
    apiSecret = "your_api_secret"
)
```
Then, we create a new client object using the NewClient function:

```go
client := Bwss.NewClient()
```
We then add a private connection for authenticated requests using the AddConnPrivate function:

```go
privateConn, err := client.AddConnPrivate(Bwss.Testnet, apiKey, apiSecret)
if err != nil {
    log.Fatal(err)
}
```
In this example, we connect to the Bybit Testnet environment, but you can also connect to the Bybit Mainnet by changing the wss.WssUrl argument to: 
```go
Bwss.Mainnet
Bwss.Testnet
```

The AddConnPrivate method takes the private WebSocket URL, API key, and API secret as arguments, and sets up a new private connection with those credentials. The AddPrivateSubs method takes a list of subscription topics (in this case, wallet and position) and a list of handler functions, which are called when a message is received on the corresponding topic.

In this example, the handler functions print out the received messages in a formatted way. The Listen method starts listening for messages on the WebSocket connection.

The AddConnPublic method takes the public WebSocket URL as an argument and sets up a new public connection. The AddPublicSubs method takes a list of subscription topics (in this case, orderbook.1.BTCUSDT) and a list of handler functions, which are called when a message is received on the corresponding topic. In this example, the handler function simply prints out the raw message.

The CloseConn method takes a connection ID as an argument and closes the corresponding WebSocket connection.

Finally, the main function sets up a signal handler to catch the interrupt signal and waits for it. When the signal is received, the wss.Close() method is called to close all open WebSocket connections.
