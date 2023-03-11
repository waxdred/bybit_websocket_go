package bybit_websocket_go

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
	stop      chan bool
	sub       int
}

func (wss *WssBybit) AddPrivateSubs(args []string, subs ...SubscribeHandler) *WssBybit {
	wss.mu.Lock()
	defer wss.mu.Unlock()
	if wss.listenner.types != "private" {
		fmt.Println("Error connection private add public Sub")
		wss.CloseConn(wss.listenner.key)
		return wss
	}
	if len(args) != len(subs) {
		wss.err = errors.New("Please add args and SubscribeHandler")
		return wss
	} else if _, ok := wss.handlePriv[wss.listenner.key]; !ok {
		return wss
	}
	if wss.err != nil {
		fmt.Println(wss.err)
		return wss
	}
	for i, arg := range args {
		switch arg {
		case "position", "execution", "order", "wallet", "greek":
			if wss.Dg {
				fmt.Println(arg, "is set")
			}
			wss = setHandler(arg, wss, &subs[i], args)
		}
	}
	return wss
}

func (wss *WssBybit) addPrivate(url Wssurl, apiKey, apiSecret string) (*WssBybit, error) {
	// Create websocket connection.
	conn, _, err := websocket.DefaultDialer.Dial(string(url), nil)
	if err != nil {
		fmt.Println("connection failed:", err)
		wss.err = err
		return wss, err
	}
	if wss.listenner == nil {
		wss.listenner = &listenner{
			key:   WssId(conn.RemoteAddr().String()),
			types: "private",
		}
	} else {
		fmt.Println("Error AddConnPrivate is already init use method listen()")
		wss.err = errors.New("Error AddConnPrivate is already init use method listen()")
		return wss, wss.err
	}
	if wss.Dg {
		fmt.Println("Connection:", conn.RemoteAddr().String())
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
		stop:      make(chan bool),
	}
	wss.handlePriv[wss.listenner.key] = &priv
	wss.nbconn += 1
	// Authenticate with API.
	wss.auth(apiKey, apiSecret)
	return wss, nil
}
