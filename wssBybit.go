package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type WssBybit struct {
	conn          *websocket.Conn
	positionBool  bool
	executionBool bool
	orderBool     bool
	walletBool    bool
	greekBool     bool
	position      func([]byte)
	execution     func([]byte)
	order         func([]byte)
	wallet        func([]byte)
	greek         func([]byte)
	messages      chan []byte
}

func (wss *WssBybit) New(url, apiKey, apiSecret string) *WssBybit {
	sign, expires := generateSignature(apiKey, apiSecret)
	// Create websocket connection.
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("connection failed:", err)
	}
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("error reading message from WebSocket:", err)
				close(wss.messages)
				return
			}
			fmt.Println(string(message))
			wss.messages <- message
		}
	}()
	// Authenticate with API.
	authPayload := map[string]interface{}{
		"op":   "auth",
		"args": []interface{}{apiKey, expires, sign},
	}
	err = conn.WriteJSON(authPayload)
	if err != nil {
		log.Fatal("authentication failed:", err)
	}
	fmt.Println("authentication successful")
	// Start goroutine to receive messages from WebSocket.
	return wss
}

func (wss *WssBybit) Subscribe(subscribe ...string) {
	for _, sub := range subscribe {
		switch sub {
		case "postion":
			wss.positionBool = true
		default:
			fmt.Println("Subscribe not found")
		}
	}
}

func (wss *WssBybit) Close() error {
	if wss.conn != nil {
		return wss.conn.Close()
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
