package blinker

import (
	"machine"
	"sync"
	"time"
)

//	type Blinker interface {
//		// Run(context.Context, *sync.WaitGroup)
//		Run(*sync.WaitGroup)
//	}
type GracefulBlinker struct {
	pwm     machine.TCC
	led     machine.Pin
	tickers *Tickers
}

func NewGracefulBlinker(pwm *machine.TCC, led machine.Pin, tickers *Tickers) *GracefulBlinker {
	led.Configure(machine.PinConfig{
		Mode: machine.PinOutput,
	})
	return &GracefulBlinker{
		pwm:     *pwm,
		led:     led,
		tickers: tickers,
	}
}

// func (g *GracefulBlinker) Run(ctx context.Context, wg *sync.WaitGroup) {
func (g *GracefulBlinker) Run(wg *sync.WaitGroup) {
	println("starting routine")
	pinUpdInterval := time.NewTicker(50 * time.Millisecond)
	defer pinUpdInterval.Stop()
	ch, err := g.pwm.Channel(g.led)
	if err != nil {
		println(err.Error())
		return
	}
	go func(wg *sync.WaitGroup, ch uint8) {
		defer wg.Done()
		for i := 0; i < 65000; i += 650 {
			<-pinUpdInterval.C
			g.pwm.Set(ch, uint32(i))
		}
		for i := 65000; i > 0; i -= 650 {
			<-pinUpdInterval.C
			g.pwm.Set(ch, uint32(i))
		}
	}(wg, ch)
}
