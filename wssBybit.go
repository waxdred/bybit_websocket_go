package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type SubscribeHandler func(*WssBybit)

type WssBybit struct {
	mu       sync.RWMutex
	priv     private
	pub      public
	nbconn   int
	nbHandle int
}

type public struct {
	conn *websocket.Conn
	recv chan []byte
	err  chan []byte
}

type private struct {
	conn      *websocket.Conn
	apiKey    string
	apiSecret string
	url       string
	position  SubscribeHandler
	execution SubscribeHandler
	order     SubscribeHandler
	wallet    SubscribeHandler
	greek     SubscribeHandler
	recv      chan []byte
	msg       []byte
	err       chan []byte
}

func (wss *WssBybit) AddPrivateSubs(topic string, args []string, subs SubscribeHandler) {
	wss.mu.Lock()
	defer wss.mu.Unlock()
	send := map[string]interface{}{
		"op":   "subscribe",
		"args": args,
	}
	switch topic {
	case "position":
		wss.priv.position = subs
		wss.priv.conn.WriteJSON(send)
		wss.nbHandle += 1
	case "execution":
		wss.priv.execution = subs
		wss.priv.conn.WriteJSON(send)
		wss.nbHandle += 1
	case "order":
		wss.priv.order = subs
		wss.priv.conn.WriteJSON(send)
		wss.nbHandle += 1
	case "wallet":
		wss.priv.wallet = subs
		wss.priv.conn.WriteJSON(send)
		wss.nbHandle += 1
	case "greek":
		wss.priv.greek = subs
		wss.priv.conn.WriteJSON(send)
		wss.nbHandle += 1
	}
}

func (wss *WssBybit) RecvPosition() *Position {
	wss.mu.Lock()
	defer wss.mu.Unlock()
	position := Position{}
	err := json.Unmarshal(wss.priv.msg, &position)
	if err != nil {
		return nil
	}
	return &position
}

func (wss *WssBybit) RecvExecution() *Execution {
	wss.mu.Lock()
	defer wss.mu.Unlock()
	execution := Execution{}
	err := json.Unmarshal(wss.priv.msg, &execution)
	if err != nil {
		return nil
	}
	return &execution
}

func (wss *WssBybit) RecvOrder() *Order {
	wss.mu.Lock()
	defer wss.mu.Unlock()
	order := Order{}
	err := json.Unmarshal(wss.priv.msg, &order)
	if err != nil {
		return nil
	}
	return &order
}

func (wss *WssBybit) RecvWallet() *Wallet {
	wss.mu.Lock()
	defer wss.mu.Unlock()
	wallet := Wallet{}
	err := json.Unmarshal(wss.priv.msg, &wallet)
	if err != nil {
		return nil
	}
	return &wallet
}

func (wss *WssBybit) RecvGreek() *Greek {
	wss.mu.Lock()
	defer wss.mu.Unlock()
	greek := Greek{}
	err := json.Unmarshal(wss.priv.msg, &greek)
	if err != nil {
		return nil
	}
	return &greek
}

func (wss *WssBybit) reconnecte() *WssBybit {
	var err error
	wss.priv.conn, _, err = websocket.DefaultDialer.Dial(wss.priv.url, nil)
	if err != nil {
		log.Fatal("connection failed:", err)
	}
	wss.auth()
	return wss
}

func (wss *WssBybit) auth() {
	sign, expires := generateSignature(apiKey, apiSecret)
	authPayload := map[string]interface{}{
		"op":   "auth",
		"args": []interface{}{wss.priv.apiKey, expires, sign},
	}
	err := wss.priv.conn.WriteJSON(authPayload)
	if err != nil {
		log.Fatal("authentication failed:", err)
	}
	fmt.Println("authentication successful")
}

func (wss *WssBybit) NewPrivate(urlPrivate, apiKey, apiSecret string) *WssBybit {
	var err error
	wss.priv.url = url
	wss.priv.apiKey = apiKey
	wss.priv.apiSecret = apiSecret
	wss.nbconn += 1
	wss.priv.recv = make(chan []byte)
	wss.priv.err = make(chan []byte)
	// Create websocket connection.
	wss.priv.conn, _, err = websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("connection failed:", err)
	}
	go func() {
		for {
			if wss.priv.conn == nil {
				wss.reconnecte()
			}
			data := make(map[string]interface{})
			_, message, err := wss.priv.conn.ReadMessage()
			if err != nil {
				log.Println("error reading message from WebSocket:", err)
				close(wss.priv.recv)
				return
			}
			json.Unmarshal(message, &data)
			if value, ok := data["success"].(bool); ok {
				if !value {
					wss.priv.err <- message
				}
			}
			// fmt.Println(string(message))
			wss.priv.recv <- message
		}
	}()
	// Authenticate with API.
	wss.auth()
	// Start goroutine to receive messages from WebSocket.
	go func() {
		ticker := time.NewTicker(time.Duration(time.Second * 20))
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				wss.mu.Lock()
				heartbeat := map[string]interface{}{
					"op":   "subscribe",
					"args": []string{"ping"},
				}
				wss.priv.conn.WriteJSON(heartbeat)
				wss.mu.Unlock()
			case err := <-wss.priv.err:
				panic(string(err))
			case message := <-wss.priv.recv:
				data := make(map[string]interface{})
				err := json.Unmarshal(message, &data)
				if err != nil {
					log.Println("error decoding message:", err)
					continue
				}
				if topic, ok := data["topic"]; ok {
					switch topic {
					case "position":
						wss.priv.position(wss)
					case "execution":
						wss.priv.execution(wss)
					case "order":
						wss.priv.order(wss)
					case "wallet":
						wss.priv.wallet(wss)
					case "greek":
						wss.priv.greek(wss)
					}
				}
			default:
			}
		}
	}()
	return wss
}

func (wss *WssBybit) Close() error {
	if wss.priv.conn != nil {
		return wss.priv.conn.Close()
	}
	return nil
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
