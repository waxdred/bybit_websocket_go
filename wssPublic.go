package bybit_websocket_go

import (
	"errors"
	"fmt"
	"strings"
	"time"

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
	stop        chan bool
}

func (wss *WssBybit) AddPublicSubs(args []string, subs ...SubscribeHandler) *WssBybit {
	wss.mu.Lock()
	defer wss.mu.Unlock()
	if wss.listenner.types != "public" {
		fmt.Println("Error connection private add public Sub")
		wss.Close()
		return wss
	}
	if _, ok := wss.handlePub[wss.listenner.key]; !ok {
		return wss
	}
	if len(args) != len(subs) {
		wss.err = errors.New("Please add args and SubscribeHandler")
		return wss
	}
	if wss.err != nil {
		return wss
	}
	for i, arg := range args {
		topic := strings.Split(arg, ".")
		if len(topic) > 0 {
			switch topic[0] {
			case "orderbook", "trade", "ticker", "kline", "liquidation",
				"kline_lt", "ticker_lt", "lt":
				wss = setHandler(topic[0], wss, &subs[i], args)
			}
		}
	}
	return wss
}

func (wss *WssBybit) addPub(url Wssurl) (*WssBybit, error) {
	// Create websocket connection.
	if wss.nbconn >= 500 {
		ticker := time.NewTicker(time.Duration(time.Minute * 1))
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if wss.nbconn < 400 {
					break
				}
			}
		}
	}

	conn, _, err := websocket.DefaultDialer.Dial(string(url), nil)
	if err != nil {
		fmt.Println("connection failed:", err)
		wss.err = err
		return wss, err
	}
	if wss.listenner == nil {
		wss.listenner = &listenner{
			key:   WssId(conn.RemoteAddr().String()),
			types: "public",
		}
	} else if wss.listenner.types != "public" {
		fmt.Println("Error fields")
		wss.err = errors.New("Error fields")
		return wss, wss.err
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
		stop:        make(chan bool),
	}
	if wss.Dg {
		fmt.Println("Connection ok")
	}
	wss.handlePub[wss.listenner.key] = &pub
	wss.nbconn += 1
	wss.reset += 1
	return wss, nil
}
