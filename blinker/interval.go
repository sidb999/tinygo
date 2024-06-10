package blinker

// import (
// 	"context"
// 	"sync"
// 	"time"
// )
//
// type Intervaler interface {
// 	Configure(time.Duration)
// 	Run(context.Context, chan<- bool)
// }
//
// type IntervalGenerator struct {
// 	interval time.Duration
// 	output   chan<- bool
// 	wg       sync.WaitGroup
// }
//
// func (i *IntervalGenerator) Configure(interval time.Duration) {
// 	i.interval = interval
// }
//
// func (i *IntervalGenerator) Run(ctx context.Context, ch chan<- bool) {
// }
