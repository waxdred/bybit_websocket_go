package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

type (
	SubscribeHandler func(*WssBybit, *SocketMessage)
)
type WssId int

type WssBybit struct {
	mu sync.RWMutex
	// priv     map[string]*private
	pub        public
	nbconn     int
	nbHandle   int
	handlePriv []*private
	handlePub  []*public
	err        error
	WssUrl     WssUrl
	Dg         bool
	run        string
}

type SocketMessage struct {
	Msg []byte
	Key string
}

func (wss *WssBybit) New(debugInfo bool) *WssBybit {
	fmt.Println(debugInfo)
	newWss := WssBybit{
		mu: sync.RWMutex{},
		// priv:     make(map[string]*private),
		nbconn:     0,
		nbHandle:   0,
		handlePriv: nil,
		handlePub:  nil,
		err:        nil,
		WssUrl:     WssUrl{},
		Dg:         true,
		run:        "",
	}
	wss = &newWss
	return wss
}

// types private | public
func (wss *WssBybit) AddConn(types string, url Wssurl, apiKey, apiSecret string) *WssBybit {
	fmt.Println(wss.Dg)
	switch types {
	case "private":
		wss.addPrivate(url, apiKey, apiSecret)
	case "public":
		wss.addPub(url)
	}
	return wss
}

func (wss *WssBybit) Listen() (WssId, error) {
	last := len(wss.handlePriv) - 1
	if wss.err != nil {
		err := wss.err
		wss.err = nil
		return -1, err
	}
	if last < 0 {
		wss.err = errors.New("Listen error not handle add")
		return -1, wss.err
	}
	if wss.run == "private" {
		if wss.handlePriv[last].sub > 0 {
			go wss.listenSocketPrivate(last)
			go wss.routineMessagePrivate(last)
		}
	} else if wss.run == "public" {
		if wss.handlePriv[last].sub > 0 {
			go wss.listenSocketPub(last)
			go wss.routineMessagePub(last)
		}
	}
	return WssId(last), nil
}

func setHandler(topic string, wss *WssBybit, subs *SubscribeHandler, args []string, idx int) *WssBybit {
	send := map[string]interface{}{
		"op":   "subscribe",
		"args": args,
	}
	switch topic {
	case "position":
		wss.handlePriv[idx].position = subs
	case "execution":
		wss.handlePriv[idx].execution = subs
	case "order":
		wss.handlePriv[idx].order = subs
	case "wallet":
		wss.handlePriv[idx].wallet = subs
	case "greek":
		wss.handlePriv[idx].greek = subs
	case "orderbook":
		wss.handlePub[idx].orderbook = subs
	case "trade":
		wss.handlePub[idx].trade = subs
	case "ticker":
		wss.handlePub[idx].ticker = subs
	case "kline":
		wss.handlePub[idx].kline = subs
	case "liquidation":
		wss.handlePub[idx].liquidation = subs
	case "kline_lt":
		wss.handlePub[idx].kline_lt = subs
	case "ticker_lt":
		wss.handlePub[idx].ticker_lt = subs
	case "lt":
		wss.handlePub[idx].lt = subs

	}
	wss.handlePriv[idx].conn.WriteJSON(send)
	wss.nbHandle += 1
	wss.handlePriv[idx].sub += 1
	return wss
}

func (sockk *SocketMessage) Unmarshal(data interface{}) error {
	if sockk != nil {
		return json.Unmarshal(sockk.Msg, &data)
	}
	return nil
}

func (wss *WssBybit) auth(idx int) {
	sign, expires := generateSignature(apiKey, apiSecret)
	authPayload := map[string]interface{}{
		"op":   "auth",
		"args": []interface{}{wss.handlePriv[idx].apiKey, expires, sign},
	}
	err := wss.handlePriv[idx].conn.WriteJSON(authPayload)
	if err != nil {
		log.Fatal("authentication failed:", err)
	}
}

func (wss *WssBybit) CloseConn(id WssId) {
	if WssId(len(wss.handlePriv)) > id && wss.handlePriv[id] != nil {
		wss.handlePriv[id].conn.Close()
	}
}

func (wss *WssBybit) Close() {
	for _, w := range wss.handlePriv {
		if w != nil && w.conn != nil {
			w.conn.Close()
		}
	}
	for _, w := range wss.handlePub {
		if w != nil && w.conn != nil {
			w.conn.Close()
		}
	}
}

func generateSignature(apiKey string, apiSecret string) (string, int64) {
	// Generate expires.
	expires := int64((time.Now().UnixNano() / int64(time.Millisecond)) + 1000)

	// Generate signature.
	signature := hmac.New(sha256.New, []byte(apiSecret))
	signature.Write([]byte(fmt.Sprintf("GET/realtime%d", expires)))
	signatureStr := fmt.Sprintf("%x", signature.Sum(nil))

	return signatureStr, expires
}
