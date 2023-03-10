package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gorilla/websocket"
)

type public struct {
	conn        *websocket.Conn
	orderbook   *SubscribeHandler
	trade       *SubscribeHandler
	ticker      *SubscribeHandler
	kline       *SubscribeHandler
	liquidation *SubscribeHandler
	kline_lt    *SubscribeHandler
	ticker_lt   *SubscribeHandler
	lt          *SubscribeHandler
	recv        chan SocketMessage
	err         chan []byte
}

func (wss *WssBybit) AddPublicSubs(args []string, subs ...SubscribeHandler) *WssBybit {
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
		topic := strings.Split(arg, ".")
		if len(topic) > 0 {
			switch topic[0] {
			case "orderbook", "trade", "ticker", "kline", "liquidation",
				"kline_lt", "ticker_lt", "lt":
				fmt.Println("send handler")
				wss = setHandler(arg, wss, &subs[i], args, last)
			}
		}
	}
	return wss
}

func (wss *WssBybit) addPub(url Wssurl) *WssBybit {
	// Create websocket connection.
	wss.run = "public"
	conn, _, err := websocket.DefaultDialer.Dial(string(url), nil)
	if err != nil {
		fmt.Println("connection failed:", err)
		wss.err = err
		return wss
	}
	pub := public{
		conn:        conn,
		orderbook:   nil,
		trade:       nil,
		ticker:      nil,
		kline_lt:    nil,
		liquidation: nil,
		kline:       nil,
		ticker_lt:   nil,
		lt:          nil,
		recv:        make(chan SocketMessage),
	}
	wss.handlePub = append(wss.handlePub, &pub)
	wss.nbconn += 1
	return wss
}
