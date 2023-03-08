package bybitwebsocket

type Position struct {
	ID           string `json:"id"`
	Topic        string `json:"topic"`
	CreationTime int64  `json:"creationTime"`
	Data         []struct {
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
