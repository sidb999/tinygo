package blinker

import (
	// "context"
	"sync"
	"time"
)

//	type Starter interface {
//		Go(wg *sync.WaitGroup)
//		Blinkers([]*GracefulBlinker) error
//		// Delay(time.Duration) error
//		// Ctx(context.Context) error
//	}
type FlawlessStarter struct {
	blinkers []*GracefulBlinker
	// delay    time.Duration
	tickers *Tickers
	// ctx     context.Context
}

func (f *FlawlessStarter) Blinkers(b []*GracefulBlinker) (err error) {
	f.blinkers = b
	return
}

// func (f *FlawlessStarter) Ctx(c context.Context) (err error) {
// 	f.ctx = c
// 	return
// }

func (f *FlawlessStarter) Go(wg *sync.WaitGroup) {
	startDelay := time.NewTicker(300 * time.Millisecond)
	defer startDelay.Stop()
	for _, blinker := range f.blinkers {
		<-startDelay.C
		wg.Add(1)
		blinker.Run(wg)
	}
}

func NewFlawlessStarter(tickers *Tickers) *FlawlessStarter {
	return &FlawlessStarter{}
}
