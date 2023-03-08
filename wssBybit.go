package bybitwebsocket

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type WssBybit struct {
	conn *websocket.Conn
}

func (wss *WssBybit) New(url, apiKey, apiSecret string) *WssBybit {
	sign, expires := generateSignature(apiKey, apiSecret)
	// Create websocket connection.
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("connection failed:", err)
	}
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
	return wss
}

func (wss *WssBybit) Close() error {
	return wss.conn.Close()
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