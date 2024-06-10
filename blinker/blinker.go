package blinker

import (
	"machine"
	"sync"
	"time"
)

type Blinker interface {
	// Run(context.Context, *sync.WaitGroup)
	Run(*sync.WaitGroup)
}
type GracefulBlinker struct {
	pwm *LockedPWM
	led machine.Pin
}

func NewGracefulBlinker(pwm *LockedPWM, led machine.Pin) Blinker {
	led.Configure(machine.PinConfig{
		Mode: machine.PinOutput,
	})
	return &GracefulBlinker{
		pwm: pwm,
		led: led,
	}
}

// func (g *GracefulBlinker) Run(ctx context.Context, wg *sync.WaitGroup) {
func (g *GracefulBlinker) Run(wg *sync.WaitGroup) {
	println("starting routine")
	ch, err := g.pwm.Get().Channel(g.led)
	if err != nil {
		println(err.Error())
		return
	}
	go func(wg *sync.WaitGroup, ch uint8) {
		defer wg.Done()
		top := g.pwm.Get().Top()
		x := top
		for {
			g.pwm.Get().Set(ch, x)
			x = x - top/100
			if x == 0 {
				x = top
			}
			time.Sleep(25 * time.Millisecond)
		}
	}(wg, ch)
}
