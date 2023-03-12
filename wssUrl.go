package bybit_websocket_go

const (
	wssSpot          = "wss://stream.bybit.com/v5/public/spot"
	wssPerpetual     = "wss://stream.bybit.com/v5/public/linear"
	wssContract      = "wss://stream.bybit.com/v5/public/inverse"
	wssOption        = "wss://stream.bybit.com/v5/public/option"
	wssSpotTest      = "wss://stream-testnet.bybit.com/v5/public/spot"
	wssPerpetualTest = "wss://stream-testnet.bybit.com/v5/public/linear"
	wssContractTest  = "wss://stream-testnet.bybit.com/v5/public/inverse"
	wssOptionTest    = "wss://stream-testnet.bybit.com/v5/public/option"
	wssPriv          = "wss://stream.bybit.com/v5/private"
	wssPrivTest      = "wss://stream-testnet.bybit.com/v5/private"
	Mainnet          = "Mainnet"
	Testnet          = "Testnet"
)

type (
	Channel string
	Wssurl  string
)

type WssUrl struct{}

func (w *Wssurl) getName() string {
	switch string(*w) {
	case "wss://stream.bybit.com/v5/public/spot":
		return "wssSpot"
	case "wss://stream.bybit.com/v5/public/linear":
		return "wssPerpetual"
	case "wss://stream.bybit.com/v5/public/inverse":
		return "wssContract"
	case "wss://stream.bybit.com/v5/public/option":
		return "wssOption"
	case "wss://stream-testnet.bybit.com/v5/public/spot":
		return "wssSpotTest"
	case "wss://stream-testnet.bybit.com/v5/public/linear":
		return "wssPerpetualTest"
	case "wss://stream-testnet.bybit.com/v5/public/inverse":
		return "wssContractTest"
	case "wss://stream-testnet.bybit.com/v5/public/option":
		return "wssOptionTest"
	case "wss://stream.bybit.com/v5/private":
		return "wssPriv"
	case "wss://stream-testnet.bybit.com/v5/private":
		return "wssPrivTest"
	default:
		return ""
	}
}

func (w *WssUrl) Spot(option Channel) Wssurl {
	switch option {
	case Mainnet:
		return wssSpot
	case Testnet:
		return wssSpotTest
	default:
		return ""
	}
}

func (w *WssUrl) Perpetual(option Channel) Wssurl {
	switch option {
	case Mainnet:
		return wssPerpetual
	case Testnet:
		return wssPerpetualTest
	default:
		return ""
	}
}

func (w *WssUrl) Contract(option Channel) Wssurl {
	switch option {
	case Mainnet:
		return wssContract
	case Testnet:
		return wssContractTest
	default:
		return ""
	}
}

func (w *WssUrl) Option(option Channel) Wssurl {
	switch option {
	case Mainnet:
		return wssOption
	case Testnet:
		return wssOptionTest
	default:
		return ""
	}
}

func (w *WssUrl) Private(option Channel) Wssurl {
	switch option {
	case Mainnet:
		return wssPriv
	case Testnet:
		return wssPrivTest
	default:
		return ""
	}
}
