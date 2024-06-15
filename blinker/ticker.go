package blinker

import "time"

type Tickers struct {
	PinChange  *time.Ticker
	StartDelay *time.Ticker
	// printTick  *time.Ticker
	StateTick *time.Ticker
}

func NewTickers() *Tickers {
	return &Tickers{
		PinChange:  time.NewTicker(50 * time.Millisecond),
		StartDelay: time.NewTicker(300 * time.Millisecond),
		// printTick:  time.NewTicker(500 * time.Millisecond),
		StateTick: time.NewTicker(6000 * time.Millisecond),
	}
}
