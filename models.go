package bybit_websocket_go

import (
	"encoding/json"
)

func pretty(data interface{}) []byte {
	prettyJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil
	}
	return prettyJSON
}

type Topic struct {
	ID           string `json:"id"`
	Topic        string `json:"topic"`
	CreationTime int64  `json:"creationTime"`
}

type Position struct {
	Topic Topic
	Data  []struct {
		PositionIdx    int    `json:"positionIdx"`
		TradeMode      int    `json:"tradeMode"`
		RiskID         int    `json:"riskId"`
		RiskLimitValue string `json:"riskLimitValue"`
		Symbol         string `json:"symbol"`
		Side           string `json:"side"`
		Size           string `json:"size"`
		EntryPrice     string `json:"entryPrice"`
		Leverage       string `json:"leverage"`
		PositionValue  string `json:"positionValue"`
		MarkPrice      string `json:"markPrice"`
		PositionIM     string `json:"positionIM"`
		PositionMM     string `json:"positionMM"`
		TakeProfit     string `json:"takeProfit"`
		StopLoss       string `json:"stopLoss"`
		TrailingStop   string `json:"trailingStop"`
		UnrealisedPnl  string `json:"unrealisedPnl"`
		CumRealisedPnl string `json:"cumRealisedPnl"`
		CreatedTime    string `json:"createdTime"`
		UpdatedTime    string `json:"updatedTime"`
		TpslMode       string `json:"tpslMode"`
		LiqPrice       string `json:"liqPrice"`
		BustPrice      string `json:"bustPrice"`
		Category       string `json:"category"`
		PositionStatus string `json:"positionStatus"`
	} `json:"data"`
}

func (p *Position) PrettyFormat() string {
	return string(pretty(p))
}

type Execution struct {
	Topic Topic
	Data  []struct {
		Category        string `json:"category"`
		Symbol          string `json:"symbol"`
		ExecFee         string `json:"execFee"`
		ExecID          string `json:"execId"`
		ExecPrice       string `json:"execPrice"`
		ExecQty         string `json:"execQty"`
		ExecType        string `json:"execType"`
		ExecValue       string `json:"execValue"`
		IsMaker         bool   `json:"isMaker"`
		FeeRate         string `json:"feeRate"`
		TradeIv         string `json:"tradeIv"`
		MarkIv          string `json:"markIv"`
		BlockTradeID    string `json:"blockTradeId"`
		MarkPrice       string `json:"markPrice"`
		IndexPrice      string `json:"indexPrice"`
		UnderlyingPrice string `json:"underlyingPrice"`
		LeavesQty       string `json:"leavesQty"`
		OrderID         string `json:"orderId"`
		OrderLinkID     string `json:"orderLinkId"`
		OrderPrice      string `json:"orderPrice"`
		OrderQty        string `json:"orderQty"`
		OrderType       string `json:"orderType"`
		StopOrderType   string `json:"stopOrderType"`
		Side            string `json:"side"`
		ExecTime        string `json:"execTime"`
		IsLeverage      string `json:"isLeverage"`
	} `json:"data"`
}

func (p *Execution) PrettyFormat() string {
	return string(pretty(p))
}

type Order struct {
	Topic Topic
	Data  []struct {
		Symbol             string `json:"symbol"`
		OrderID            string `json:"orderId"`
		Side               string `json:"side"`
		OrderType          string `json:"orderType"`
		CancelType         string `json:"cancelType"`
		Price              string `json:"price"`
		Qty                string `json:"qty"`
		OrderIv            string `json:"orderIv"`
		TimeInForce        string `json:"timeInForce"`
		OrderStatus        string `json:"orderStatus"`
		OrderLinkID        string `json:"orderLinkId"`
		LastPriceOnCreated string `json:"lastPriceOnCreated"`
		ReduceOnly         bool   `json:"reduceOnly"`
		LeavesQty          string `json:"leavesQty"`
		LeavesValue        string `json:"leavesValue"`
		CumExecQty         string `json:"cumExecQty"`
		CumExecValue       string `json:"cumExecValue"`
		AvgPrice           string `json:"avgPrice"`
		BlockTradeID       string `json:"blockTradeId"`
		PositionIdx        int    `json:"positionIdx"`
		CumExecFee         string `json:"cumExecFee"`
		CreatedTime        string `json:"createdTime"`
		UpdatedTime        string `json:"updatedTime"`
		RejectReason       string `json:"rejectReason"`
		StopOrderType      string `json:"stopOrderType"`
		TriggerPrice       string `json:"triggerPrice"`
		TakeProfit         string `json:"takeProfit"`
		StopLoss           string `json:"stopLoss"`
		TpTriggerBy        string `json:"tpTriggerBy"`
		SlTriggerBy        string `json:"slTriggerBy"`
		TriggerDirection   int    `json:"triggerDirection"`
		TriggerBy          string `json:"triggerBy"`
		CloseOnTrigger     bool   `json:"closeOnTrigger"`
		Category           string `json:"category"`
	} `json:"data"`
}

func (p *Order) PrettyFormat() string {
	return string(pretty(p))
}

type Wallet struct {
	Topic Topic
	Data  []struct {
		AccountIMRate          string `json:"accountIMRate"`
		AccountMMRate          string `json:"accountMMRate"`
		TotalEquity            string `json:"totalEquity"`
		TotalWalletBalance     string `json:"totalWalletBalance"`
		TotalMarginBalance     string `json:"totalMarginBalance"`
		TotalAvailableBalance  string `json:"totalAvailableBalance"`
		TotalPerpUPL           string `json:"totalPerpUPL"`
		TotalInitialMargin     string `json:"totalInitialMargin"`
		TotalMaintenanceMargin string `json:"totalMaintenanceMargin"`
		Coin                   []struct {
			Coin                string `json:"coin"`
			Equity              string `json:"equity"`
			UsdValue            string `json:"usdValue"`
			WalletBalance       string `json:"walletBalance"`
			AvailableToWithdraw string `json:"availableToWithdraw"`
			AvailableToBorrow   string `json:"availableToBorrow"`
			BorrowAmount        string `json:"borrowAmount"`
			AccruedInterest     string `json:"accruedInterest"`
			TotalOrderIM        string `json:"totalOrderIM"`
			TotalPositionIM     string `json:"totalPositionIM"`
			TotalPositionMM     string `json:"totalPositionMM"`
			UnrealisedPnl       string `json:"unrealisedPnl"`
			CumRealisedPnl      string `json:"cumRealisedPnl"`
			Bonus               string `json:"bonus"`
		} `json:"coin"`
		AccountType string `json:"accountType"`
	} `json:"data"`
}

func (p *Wallet) PrettyFormat() string {
	return string(pretty(p))
}

type Greek struct {
	Topic Topic
	Data  []struct {
		BaseCoin   string `json:"baseCoin"`
		TotalDelta string `json:"totalDelta"`
		TotalGamma string `json:"totalGamma"`
		TotalVega  string `json:"totalVega"`
		TotalTheta string `json:"totalTheta"`
	} `json:"data"`
}

func (p *Greek) PrettyFormat() string {
	return string(pretty(p))
}
