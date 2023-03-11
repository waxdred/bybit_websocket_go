package bybit_websocket_go

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
)

type (
	SubscribeHandler func(*WssBybit, *SocketMessage)
)
type WssId string

type WssBybit struct {
	mu sync.RWMutex
	// priv     map[string]*private
	pub        public
	nbconn     int
	nbHandle   int
	handlePriv map[WssId]*private
	handlePub  map[WssId]*public
	err        error
	WssUrl     WssUrl
	Dg         bool
	run        string
	listenner  *listenner
}

type listenner struct {
	key   WssId
	types string
}

type SocketMessage struct {
	Msg []byte
	Key string
}

func (wss *WssBybit) New(debugInfo bool) *WssBybit {
	wss = &WssBybit{
		mu:         sync.RWMutex{},
		nbconn:     0,
		nbHandle:   0,
		handlePriv: make(map[WssId]*private),
		handlePub:  make(map[WssId]*public),
		err:        nil,
		WssUrl:     WssUrl{},
		Dg:         debugInfo,
		run:        "",
		listenner:  nil,
	}
	return wss
}

// types private | public
func (wss *WssBybit) AddConnPrivate(url Wssurl, apiKey, apiSecret string) *WssBybit {
	wss.addPrivate(url, apiKey, apiSecret)
	return wss
}

func (wss *WssBybit) AddConnPublic(url Wssurl) *WssBybit {
	wss.addPub(url)
	return wss
}

func (wss *WssBybit) Listen() (WssId, error) {
	key := wss.listenner.key
	if wss.err != nil || wss.listenner == nil {
		err := wss.err
		fmt.Println(err)
		wss.err = nil
		return WssId(""), err
	}
	if wss.listenner.types == "private" {
		go wss.listenSocketPrivate(wss.listenner.key)
		go wss.routineMessagePrivate(wss.listenner.key)
	} else if wss.listenner.types == "public" {
		go wss.listenSocketPub(wss.listenner.key)
		go wss.routineMessagePub(wss.listenner.key)
	}
	wss.listenner = nil
	return key, nil
}

func setHandler(topic string, wss *WssBybit, subs *SubscribeHandler, args []string) *WssBybit {
	if wss.listenner == nil {
		return wss
	}
	send := map[string]interface{}{
		"op":   "subscribe",
		"args": args,
	}
	switch topic {
	case "position":
		wss.handlePriv[wss.listenner.key].position = subs
		wss.handlePriv[wss.listenner.key].conn.WriteJSON(send)
		wss.handlePriv[wss.listenner.key].sub += 1
	case "execution":
		wss.handlePriv[wss.listenner.key].execution = subs
		wss.handlePriv[wss.listenner.key].conn.WriteJSON(send)
		wss.handlePriv[wss.listenner.key].sub += 1
	case "order":
		wss.handlePriv[wss.listenner.key].order = subs
		wss.handlePriv[wss.listenner.key].conn.WriteJSON(send)
		wss.handlePriv[wss.listenner.key].sub += 1
	case "wallet":
		wss.handlePriv[wss.listenner.key].wallet = subs
		wss.handlePriv[wss.listenner.key].conn.WriteJSON(send)
		wss.handlePriv[wss.listenner.key].sub += 1
	case "greek":
		wss.handlePriv[wss.listenner.key].greek = subs
		wss.handlePriv[wss.listenner.key].conn.WriteJSON(send)
		wss.handlePriv[wss.listenner.key].sub += 1
	case "orderbook":
		if wss.handlePub[wss.listenner.key] != nil {
			wss.handlePub[wss.listenner.key].orderbook = subs
			wss.handlePub[wss.listenner.key].conn.WriteJSON(send)
		}
	case "trade":
		if wss.handlePub[wss.listenner.key] != nil {
			wss.handlePub[wss.listenner.key].trade = subs
			wss.handlePub[wss.listenner.key].conn.WriteJSON(send)
		}
	case "ticker":
		if wss.handlePub[wss.listenner.key] != nil {
			wss.handlePub[wss.listenner.key].ticker = subs
			wss.handlePub[wss.listenner.key].conn.WriteJSON(send)
		}
	case "kline":
		if wss.handlePub[wss.listenner.key] != nil {
			wss.handlePub[wss.listenner.key].kline = subs
			wss.handlePub[wss.listenner.key].conn.WriteJSON(send)
		}
	case "liquidation":
		if wss.handlePub[wss.listenner.key] != nil {
			wss.handlePub[wss.listenner.key].liquidation = subs
			wss.handlePub[wss.listenner.key].conn.WriteJSON(send)
		}
	case "kline_lt":
		if wss.handlePub[wss.listenner.key] != nil {
			wss.handlePub[wss.listenner.key].kline_lt = subs
			wss.handlePub[wss.listenner.key].conn.WriteJSON(send)
		}
	case "ticker_lt":
		if wss.handlePub[wss.listenner.key] != nil {
			wss.handlePub[wss.listenner.key].ticker_lt = subs
			wss.handlePub[wss.listenner.key].conn.WriteJSON(send)
		}
	case "lt":
		if wss.handlePub[wss.listenner.key] != nil {
			wss.handlePub[wss.listenner.key].lt = subs
			wss.handlePub[wss.listenner.key].conn.WriteJSON(send)
		}

	}
	wss.nbHandle += 1
	return wss
}

func (sockk *SocketMessage) Unmarshal(data interface{}) error {
	if sockk != nil {
		return json.Unmarshal(sockk.Msg, &data)
	}
	return nil
}

func (wss *WssBybit) auth(apiKey, apiSecret string) {
	sign, expires := generateSignature(apiKey, apiSecret)
	authPayload := map[string]interface{}{
		"op":   "auth",
		"args": []interface{}{wss.handlePriv[wss.listenner.key].apiKey, expires, sign},
	}
	err := wss.handlePriv[wss.listenner.key].conn.WriteJSON(authPayload)
	if err != nil {
		log.Fatal("authentication failed:", err)
	}
}

func (wss *WssBybit) CloseConn(id WssId) {
	fmt.Println("Close connection:", id)
	if w, ok := wss.handlePriv[id]; ok {
		w.stop <- true
		w.conn.Close()
		delete(wss.handlePriv, id)
	} else if w, ok := wss.handlePub[id]; ok {
		w.stop <- true
		w.conn.Close()
		delete(wss.handlePub, id)
	}
}

func (wss *WssBybit) Close() {
	for _, w := range wss.handlePriv {
		if w != nil && w.conn != nil {
			fmt.Println("close conn")
			w.stop <- true
			w.conn.Close()
		}
	}
	for _, w := range wss.handlePub {
		if w != nil && w.conn != nil {
			fmt.Println("close conn")
			w.stop <- true
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
