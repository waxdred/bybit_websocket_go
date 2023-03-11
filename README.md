# Bybit WebSocket Go Library
![GitHub Workflow Status](https://github.com/waxdred/bybit_websocket_go/actions/workflows/go.yml/badge.svg)
![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)
### This is a Go package for connecting to Bybit's WebSocket API. It provides a simple way to connect to both public and private streams, and includes support for subscribing to and processing real-time market data, trading information, and account information.

## Table of Contents
- [Installation](#Installation)
- [Usage](#Usage)

# Installation
To use this package, you need to have Go installed on your system. You can then install the package by running:

```go
go get github.com/waxdred/bybit_websocket_go
```
# Usage
To use the package, first import it:

```go
import Bwss "github.com/waxdred/bybit_websocket_go"
```
## Connection
To create a new connection to the Bybit WebSocket API, use the New function of the WssBybit struct:

```go
wss := new(Bwss.WssBybit).New(true)
```
You can pass a boolean parameter to enable or disable debugging messages.

To close all open connections, you can use the Close function:

```go
defer wss.Close()
```
## Public streams
To subscribe to a public stream, use the AddConnPublic function of the WssBybit struct, which takes the URL of the stream as its parameter:

```go

id, _ := wss.AddConnPublic(wss.WssUrl.Perpetual(Bwss.Mainnet))
```
The function returns an ID that can be used to close the connection later.

To subscribe to a channel on the stream, use the AddPublicSubs function, which takes a list of channel names and a handler function as its parameters:

```go
wss.AddPublicSubs([]string{"orderBookL2_25.BTCUSD"}, func(wss *Bwss.WssBybit, sockk *Bwss.SocketMessage) {
    fmt.Println(string(sockk.Msg))
})
```
## Private streams
To subscribe to a private stream, use the AddConnPrivate function of the WssBybit struct, which takes the URL of the stream, your API key, and your API secret as its parameters:

```go
wss.AddConnPrivate(wss.WssUrl.Private(Bwss.Testnet), apiKey, apiSecret)
```
To subscribe to a channel on the stream, use the AddPrivateSubs function, which takes a list of channel names and a handler function for each channel as its parameters:

```go
wss.AddPrivateSubs([]string{"wallet", "position"}, func(wss *Bwss.WssBybit, sockk *Bwss.SocketMessage) {
    wallet := Bwss.Wallet{}
    sockk.Unmarshal(&wallet)
    log.Println(wallet.PrettyFormat())
}, func(wss *Bwss.WssBybit, sockk *Bwss.SocketMessage) {
    position := Bwss.Position{}
    sockk.Unmarshal(&position)
    log.Println(position.PrettyFormat())
})
```

## Testnet vs Mainnet
Bybit offers two separate environments for testing and live trading: Testnet and Mainnet. The Testnet environment allows you to test your strategies and code without risking real funds, while the Mainnet environment is for actual trading.

In the code, you can select which environment to use by passing either Bwss.Testnet or Bwss.Mainnet as the parameter to the WssUrl functions:

``go
wss.WssUrl.Perpetual(Bwss.Testnet) // for Testnet
wss.WssUrl.Perpetual(Bwss
```
Testnet and Mainnet are two separate environments provided by Bybit for testing and live trading, respectively. Testnet is a sandbox environment that allows you to test your code and strategies without risking real funds, while Mainnet is the production environment used for actual trading.

In the provided code, you can select which environment to use by passing either Bwss.Testnet or Bwss.Mainnet as a parameter to the WssUrl functions. For example, to connect to the Testnet perpetual contract endpoint, you can use:

```go
wss.AddConnPublic(wss.WssUrl.Perpetual(Bwss.Testnet)).
```
And to connect to the Mainnet spot trading endpoint, you can use:

```go
wss.AddConnPublic(wss.WssUrl.Spot(Bwss.Mainnet)).
````
It's important to note that different API keys are required for the Testnet and Mainnet environments. Be sure to generate separate API keys for each environment and use them accordingly.

For more examples and detailed usage instructions, please see the GoDoc [documentation](https://pkg.go.dev/github.com/waxdred/bybit_websocket_go?utm_source=godoc).

## Contributing
If you find a bug or would like to contribute to this library, please open an issue or pull request on the GitHub repository.

## License
This library is licensed under the MIT License. See the LICENSE file for details.
