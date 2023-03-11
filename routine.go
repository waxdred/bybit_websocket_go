package bybit_websocket_go

import (
	"encoding/json"
	"fmt"
	"time"
)

func (wss *WssBybit) listenSocketPrivate(key WssId) {
	for {
		data := make(map[string]interface{})
		if _, ok := wss.handlePriv[key]; !ok {
			fmt.Println("close listen socket")
			return
		}
		_, message, err := wss.handlePriv[key].conn.ReadMessage()
		if err != nil {
			fmt.Println("error reading message from WebSocket:", err)
			delete(wss.handlePriv, key)
			return
		}
		json.Unmarshal(message, &data)
		if value, ok := data["op"].(string); ok {
			if value == "pong" && wss.Dg {
				fmt.Println(string(message))
			}
		}
		wss.handlePriv[key].recv <- SocketMessage{
			Msg: message,
			Key: wss.handlePriv[key].conn.RemoteAddr().String(),
		}
	}
}

func (wss *WssBybit) listenSocketPub(key WssId) {
	for {
		data := make(map[string]interface{})
		if _, ok := wss.handlePub[key]; !ok {
			fmt.Println("close listen socket")
			return
		}
		_, message, err := wss.handlePub[key].conn.ReadMessage()
		if err != nil {
			fmt.Println("error reading message from WebSocket:", err)
			if wss.handlePub[key] != nil {
				delete(wss.handlePub, key)
			}
			return
		}
		// fmt.Println(string(message))
		if wss.Dg {
			fmt.Println(string(message))
		}
		json.Unmarshal(message, &data)
		if _, ok := data["success"].(bool); ok {
			if wss.Dg {
				fmt.Println(data)
			}
		}
		if value, ok := data["op"].(string); ok {
			if value == "pong" && wss.Dg {
				fmt.Println(string(message))
			}
		}
		wss.handlePub[key].recv <- SocketMessage{
			Msg: message,
			Key: wss.handlePub[key].conn.RemoteAddr().String(),
		}
	}
}

func (wss *WssBybit) routineMessagePub(key WssId) {
	ticker := time.NewTicker(time.Duration(time.Second * 20))
	defer ticker.Stop()
	for {
		select {
		case stop := <-wss.handlePub[key].stop:
			if stop {
				close(wss.handlePub[key].stop)
				close(wss.handlePub[key].recv)
				return
			}
		case <-ticker.C:
			heartbeat := map[string]interface{}{
				"op":   "ping",
				"args": []string{"1675418560633"},
			}
			wss.handlePub[key].conn.WriteJSON(heartbeat)
		case sockk := <-wss.handlePub[key].recv:
			data := make(map[string]interface{})
			err := json.Unmarshal(sockk.Msg, &data)
			if err != nil {
				fmt.Println("error decoding message:", err)
				close(wss.handlePub[key].recv)
				close(wss.handlePub[key].stop)
				wss.CloseConn(key)
				return
			}
			if wss.Dg {
				fmt.Println(string(sockk.Msg))
			}
			if topic, ok := data["topic"]; ok {
				switch topic {
				case "orderbook":
					if wss.handlePub[key].orderbook != nil {
						(*wss.handlePub[key].orderbook)(wss, &sockk)
					}
				case "trade":
					if wss.handlePub[key].trade != nil {
						(*wss.handlePub[key].trade)(wss, &sockk)
					}
				case "ticker":
					if wss.handlePub[key].ticker != nil {
						(*wss.handlePub[key].ticker)(wss, &sockk)
					}
				case "kline":
					if wss.handlePub[key].kline != nil {
						(*wss.handlePub[key].kline)(wss, &sockk)
					}
				case "liquidation":
					if wss.handlePub[key].liquidation != nil {
						(*wss.handlePub[key].liquidation)(wss, &sockk)
					}
				case "kline_lt":
					if wss.handlePub[key].kline_lt != nil {
						(*wss.handlePub[key].kline_lt)(wss, &sockk)
					}
				case "tickers_lt":
					if wss.handlePub[key].ticker_lt != nil {
						(*wss.handlePub[key].kline)(wss, &sockk)
					}
				case "lt":
					if wss.handlePub[key].lt != nil {
						(*wss.handlePub[key].lt)(wss, &sockk)
					}
				}
			}
		default:
		}
	}
}

func (wss *WssBybit) routineMessagePrivate(key WssId) {
	ticker := time.NewTicker(time.Duration(time.Second * 20))
	defer ticker.Stop()
	for {
		select {
		case <-wss.handlePriv[key].stop:
			close(wss.handlePriv[key].stop)
			close(wss.handlePriv[key].recv)
			return
		case <-ticker.C:
			heartbeat := map[string]interface{}{
				"op":   "ping",
				"args": []string{"1675418560633"},
			}
			wss.handlePriv[key].conn.WriteJSON(heartbeat)
		case sockk := <-wss.handlePriv[key].recv:
			data := make(map[string]interface{})
			err := json.Unmarshal(sockk.Msg, &data)
			if err != nil {
				fmt.Println("error decoding message:", err)
				return
			}
			if wss.Dg {
				fmt.Println(string(sockk.Msg))
			}
			if topic, ok := data["topic"]; ok {
				switch topic {
				case "position":
					if wss.handlePriv[key].position != nil {
						(*wss.handlePriv[key].position)(wss, &sockk)
					}
				case "execution":
					if wss.handlePriv[key].execution != nil {
						(*wss.handlePriv[key].execution)(wss, &sockk)
					}
				case "order":
					if wss.handlePriv[key].order != nil {
						(*wss.handlePriv[key].order)(wss, &sockk)
					}
				case "wallet":
					if wss.handlePriv[key].wallet != nil {
						(*wss.handlePriv[key].wallet)(wss, &sockk)
					}
				case "greek":
					if wss.handlePriv[key].greek != nil {
						(*wss.handlePriv[key].greek)(wss, &sockk)
					}
				}
			}
		default:
		}
	}
}
