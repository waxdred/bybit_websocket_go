package main

import (
	"encoding/json"
	"fmt"
	"time"
)

func (wss *WssBybit) listenSocketPrivate(idx int) {
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
			fmt.Println("error reading message from WebSocket:", err)
			if wss.handlePriv[idx] != nil {
				if wss.handlePriv[idx].recv != nil {
					close(wss.handlePriv[idx].recv)
				}
				if wss.handlePriv[idx].err != nil {
					close(wss.handlePriv[idx].err)
				}
			}
			wss.handlePriv[idx] = nil
			return
		}
		fmt.Println(string(message), wss.Dg)
		if wss.Dg {
			fmt.Println(string(message))
		}
		json.Unmarshal(message, &data)
		if value, ok := data["success"].(bool); ok {
			if !value {
				wss.handlePriv[idx].err <- message
			}
		}
		if value, ok := data["op"].(string); ok {
			if value == "pong" && wss.Dg {
				fmt.Println(message)
			}
		}
		wss.handlePriv[idx].recv <- SocketMessage{
			Msg: message,
			Key: wss.handlePriv[idx].conn.RemoteAddr().String(),
		}
	}
}

func (wss *WssBybit) listenSocketPub(idx int) {
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
			fmt.Println("error reading message from WebSocket:", err)
			if wss.handlePriv[idx] != nil {
				if wss.handlePriv[idx].recv != nil {
					close(wss.handlePriv[idx].recv)
				}
				if wss.handlePriv[idx].err != nil {
					close(wss.handlePriv[idx].err)
				}
			}
			wss.handlePriv[idx] = nil
			return
		}
		fmt.Println(string(message), wss.Dg)
		if wss.Dg {
			fmt.Println(string(message))
		}
		json.Unmarshal(message, &data)
		if value, ok := data["success"].(bool); ok {
			if !value {
				wss.handlePriv[idx].err <- message
			}
		}
		if value, ok := data["op"].(string); ok {
			if value == "pong" && wss.Dg {
				fmt.Println(message)
			}
		}
		wss.handlePriv[idx].recv <- SocketMessage{
			Msg: message,
			Key: wss.handlePriv[idx].conn.RemoteAddr().String(),
		}
	}
}

func (wss *WssBybit) routineMessagePub(idx int) {
	ticker := time.NewTicker(time.Duration(time.Second * 20))
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			heartbeat := map[string]interface{}{
				"op":   "ping",
				"args": []string{"1675418560633"},
			}
			wss.handlePriv[idx].conn.WriteJSON(heartbeat)
		case err := <-wss.handlePub[idx].err:
			fmt.Println(err)
			return
		case sockk := <-wss.handlePub[idx].recv:
			data := make(map[string]interface{})
			err := json.Unmarshal(sockk.Msg, &data)
			if err != nil {
				fmt.Println("error decoding message:", err)
				continue
			}
			if topic, ok := data["topic"]; ok {
				switch topic {
				case "orderbook":
					if wss.handlePub[idx].orderbook != nil {
						(*wss.handlePub[idx].orderbook)(wss, &sockk)
					}
				case "trade":
					if wss.handlePub[idx].trade != nil {
						(*wss.handlePub[idx].trade)(wss, &sockk)
					}
				case "ticker":
					if wss.handlePub[idx].ticker != nil {
						(*wss.handlePub[idx].ticker)(wss, &sockk)
					}
				case "kline":
					if wss.handlePub[idx].kline != nil {
						(*wss.handlePub[idx].kline)(wss, &sockk)
					}
				case "liquidation":
					if wss.handlePub[idx].liquidation != nil {
						(*wss.handlePub[idx].liquidation)(wss, &sockk)
					}
				case "kline_lt":
					if wss.handlePub[idx].kline_lt != nil {
						(*wss.handlePub[idx].kline_lt)(wss, &sockk)
					}
				case "tickers_lt":
					if wss.handlePub[idx].ticker_lt != nil {
						(*wss.handlePub[idx].kline)(wss, &sockk)
					}
				case "lt":
					if wss.handlePub[idx].lt != nil {
						(*wss.handlePub[idx].lt)(wss, &sockk)
					}
				}
			}
		default:
		}
	}
}

func (wss *WssBybit) routineMessagePrivate(idx int) {
	ticker := time.NewTicker(time.Duration(time.Second * 20))
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			heartbeat := map[string]interface{}{
				"op":   "ping",
				"args": []string{"1675418560633"},
			}
			wss.handlePriv[idx].conn.WriteJSON(heartbeat)
		case err := <-wss.handlePriv[idx].err:
			fmt.Println(err)
			return
		case sockk := <-wss.handlePriv[idx].recv:
			data := make(map[string]interface{})
			err := json.Unmarshal(sockk.Msg, &data)
			if err != nil {
				fmt.Println("error decoding message:", err)
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
