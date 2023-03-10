package main

import (
	"errors"
	"fmt"

	"github.com/gorilla/websocket"
)

type private struct {
	conn      *websocket.Conn
	apiKey    string
	apiSecret string
	position  *SubscribeHandler
	execution *SubscribeHandler
	order     *SubscribeHandler
	wallet    *SubscribeHandler
	greek     *SubscribeHandler
	recv      chan SocketMessage
	msg       []byte
	err       chan []byte
	sub       int
}

func (wss *WssBybit) AddPrivateSubs(args []string, subs ...SubscribeHandler) *WssBybit {
	wss.mu.Lock()
	defer wss.mu.Unlock()
	if len(args) != len(subs) {
		wss.err = errors.New("Please add args and SubscribeHandler")
		return wss
	}
	if wss.err != nil {
		return wss
	}
	last := len(wss.handlePriv) - 1
	if last < 0 {
		wss.err = errors.New("AddPrivateSubs error Please AddConn")
		return wss
	}
	for i, arg := range args {
		switch arg {
		case "position", "execution", "order", "wallet", "greek":
			fmt.Println("send handler")
			wss = setHandler(arg, wss, &subs[i], args, last)
		}
	}
	return wss
}

func (wss *WssBybit) addPrivate(url Wssurl, apiKey, apiSecret string) *WssBybit {
	// Create websocket connection.
	wss.run = "private"
	conn, _, err := websocket.DefaultDialer.Dial(string(url), nil)
	if err != nil {
		fmt.Println("connection failed:", err)
		wss.err = err
		return wss
	}
	priv := private{
		conn:      conn,
		apiKey:    apiKey,
		apiSecret: apiSecret,
		position:  nil,
		execution: nil,
		order:     nil,
		wallet:    nil,
		greek:     nil,
		recv:      make(chan SocketMessage),
		sub:       0,
	}
	wss.handlePriv = append(wss.handlePriv, &priv)
	wss.nbconn += 1
	// Authenticate with API.
	wss.auth(len(wss.handlePriv) - 1)
	return wss
}
