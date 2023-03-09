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

	"github.com/gorilla/websocket"
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
}

type public struct {
	conn *websocket.Conn
	recv chan []byte
	err  chan []byte
}

type SocketMessage struct {
	Msg []byte
	Key string
}

type private struct {
	conn      *websocket.Conn
	listen    bool
	apiKey    string
	apiSecret string
	url       string
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

func (wss *WssBybit) New() *WssBybit {
	newWss := WssBybit{
		mu: sync.RWMutex{},
		// priv:     make(map[string]*private),
		nbconn:     0,
		nbHandle:   0,
		handlePriv: nil,
		handlePub:  nil,
		err:        nil,
	}
	wss = &newWss
	return wss
}

// types private | public
func (wss *WssBybit) AddConn(types, url, apiKey, apiSecret string) *WssBybit {
	switch types {
	case "private":
		wss.addPrivate(url, apiKey, apiSecret)
	case "public":
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
	if wss.handlePriv[last].sub > 0 {
		go wss.listenSocket(last)
		go wss.RoutineMessagePrivate(last)
	}
	return WssId(last), nil
}

func (wss *WssBybit) AddPrivateSubs(topic string, args []string, subs SubscribeHandler) *WssBybit {
	wss.mu.Lock()
	defer wss.mu.Unlock()
	last := len(wss.handlePriv) - 1
	if last < 0 {
		wss.err = errors.New("AddPrivateSubs error Please AddConn")
		return wss
	}
	switch topic {
	case "position", "execution", "order", "wallet", "greek":
		fmt.Println("send handler")
		wss = setHandler(topic, wss, &subs, args, last)
	}
	return wss
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
	fmt.Println("authentication successful")
}

func (wss *WssBybit) listenSocket(idx int) {
	for {
		data := make(map[string]interface{})
		if wss.handlePriv[idx] == nil {
			fmt.Println("close listen socket")
			close(wss.handlePriv[idx].recv)
			close(wss.handlePriv[idx].err)
			return
		}
		_, message, err := wss.handlePriv[idx].conn.ReadMessage()
		if err != nil {
			log.Println("error reading message from WebSocket:", err)
			close(wss.handlePriv[idx].recv)
			close(wss.handlePriv[idx].err)
			wss.handlePriv[idx] = nil
			return
		}
		// fmt.Println(string(message))
		json.Unmarshal(message, &data)
		if value, ok := data["success"].(bool); ok {
			if !value {
				wss.handlePriv[idx].err <- message
			}
		}
		wss.handlePriv[idx].recv <- SocketMessage{
			Msg: message,
			Key: wss.handlePriv[idx].conn.RemoteAddr().String(),
		}
	}
}

func (wss *WssBybit) addPrivate(urlPrivate, apiKey, apiSecret string) *WssBybit {
	// Create websocket connection.
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Println("connection failed:", err)
		wss.err = err
		return wss
	}
	priv := private{
		conn:      conn,
		listen:    false,
		apiKey:    apiKey,
		apiSecret: apiSecret,
		url:       url,
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

func (wss *WssBybit) RoutineMessagePrivate(idx int) {
	ticker := time.NewTicker(time.Duration(time.Second * 20))
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// log.Println("pong")
			heartbeat := map[string]interface{}{
				"op":   "subscribe",
				"args": []string{"ping"},
			}
			wss.handlePriv[idx].conn.WriteJSON(heartbeat)
		case err := <-wss.handlePriv[idx].err:
			panic(string(err))
		case sockk := <-wss.handlePriv[idx].recv:
			data := make(map[string]interface{})
			err := json.Unmarshal(sockk.Msg, &data)
			if err != nil {
				log.Println("error decoding message:", err)
				continue
			}
			if topic, ok := data["topic"]; ok {
				switch topic {
				case "position":
					if wss.handlePriv[idx].position != nil {
						(*wss.handlePriv[idx].position)(wss, &sockk)
					}
				case "execution":
					if wss.handlePriv[idx].execution != nil {
						(*wss.handlePriv[idx].execution)(wss, &sockk)
					}
				case "order":
					if wss.handlePriv[idx].order != nil {
						(*wss.handlePriv[idx].order)(wss, &sockk)
					}
				case "wallet":
					if wss.handlePriv[idx].wallet != nil {
						(*wss.handlePriv[idx].wallet)(wss, &sockk)
					}
				case "greek":
					if wss.handlePriv[idx].greek != nil {
						(*wss.handlePriv[idx].greek)(wss, &sockk)
					}
				}
			}
		default:
		}
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
