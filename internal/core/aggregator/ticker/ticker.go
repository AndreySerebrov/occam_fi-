package ticker

import "time"

/*
	PriceStreamSubscriber interface
*/

type Ticker string

type TickerPrice struct {
	Ticker Ticker
	Time   time.Time
	Price  string // decimal value. example: "0", "10", "12.2", "13.2345122"
}

// PriceStreamSubscriber
type PriceStreamSubscriber interface {
	// I didn't get the purpose of arguments for this method
	// So, my suggestion is:
	SubscribePriceStream() (chan TickerPrice, chan error)
}
