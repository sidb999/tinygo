package blinker

import (
	"context"
	"sync"
	"time"
)

type Starter interface {
	Go(wg *sync.WaitGroup)
	Blinkers([]Blinker) error
	Delay(time.Duration) error
	// Ctx(context.Context) error
}
type FlawlessStarter struct {
	blinkers []Blinker
	delay    time.Duration
	ctx      context.Context
}

func (f *FlawlessStarter) Blinkers(b []Blinker) (err error) {
	f.blinkers = b
	return
}

func (f *FlawlessStarter) Delay(d time.Duration) (err error) {
	f.delay = d
	return
}

// func (f *FlawlessStarter) Ctx(c context.Context) (err error) {
// 	f.ctx = c
// 	return
// }

func (f *FlawlessStarter) Go(wg *sync.WaitGroup) {
	for _, blinker := range f.blinkers {
		wg.Add(1)
		// blinker.Run(f.ctx, wg)
		blinker.Run(wg)
	}
}

func NewFlawlessStarter() Starter {
	return &FlawlessStarter{}
}
