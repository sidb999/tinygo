package main

import (
	"context"
	"machine"
	"sync"
	"time"
)

// diginal pins { 3, 5, 6, 9, 10, 11 };
// Pins that works with PWM LEDs:
// 5,6

//	type Starter interface {
//		Go()
//	}
//
//	type Blinker interface {
//		run(context.Context, *sync.WaitGroup) error
//	}
//
// type FlawlessStarted struct{}
//
//	type GracefulBlinker struct {
//		led machine.Pin
//	}
//
//	func NewGracefulBlinker(pin machine.Pin) *GracefulBlinker {
//		return &GracefulBlinker{
//			led: pin,
//		}
//	}
//
//	func (g *GracefulBlinker) run(ctx context.Context, wg *sync.WaitGroup) (err error) {
//		g.led.Configure(machine.PinConfig{
//			Mode: machine.PinOutput,
//		})
//		go func(ctx context.Context, wg *sync.WaitGroup) {
//			defer wg.Done()
//			// blinkInterval := time.NewTicker(1 * time.Second)
//
//			// select {
//			// case <-blinkInterval.C:
//			// 	g.led.Low
//			// }
//		}(ctx, wg)
//
//		return
//	}
func listenToSignals(cancel context.CancelFunc) {
	c := make(chan int, 1)
	<-c
	cancel()
}

func waitForDurationShutdown(cancel context.CancelFunc, dur time.Duration) {
	<-time.After(dur)
	cancel()
}

func main() {
	var wg sync.WaitGroup
	logCh := make(chan uint32, 1)
	// duration, _ := time.ParseDuration("1h")
	wg.Add(3)
	_, cancel := context.WithCancel(context.Background())
	go listenToSignals(cancel)
	// go waitForDurationShutdown(cancel, duration)

	// blinkInterval := time.NewTicker(5 * time.Millisecond)
	// defer blinkInterval.Stop()

	pin := machine.D2
	pin.Configure(machine.PinConfig{
		Mode: machine.PinOutput,
	})

	var pwmCtl machine.PWM
	var period uint64
	var ch uint8
	var err error

	err = pwmCtl.Configure(machine.PWMConfig{})
	if err != nil {
		println(err.Error())
		return
	}
	period = pwmCtl.Period()
	println(period)
	// blinkInterval := time.NewTicker(10 * time.Duration(period))
	// defer blinkInterval.Stop()

	ch, err = pwmCtl.Channel(machine.PD5)
	if err != nil {
		println("D5 (Timer0): ", err.Error())
		// return
	} else {
		println("D5 (Timer0): OK")
	}
	ch, err = pwmCtl.Channel(machine.PD6)
	if err != nil {
		println("D6: (Timer0)", err.Error())
		// return
	} else {
		println("D6 (Timer0): OK")
	}

	top := pwmCtl.Top()
	x := top
	go func() {
		defer wg.Done()
		for {
			// <-blinkInterval.C
			logCh <- x
			pwmCtl.Set(ch, x)
			x = x - top/100
			if x == 0 {
				x = top
			}
			time.Sleep(10 * time.Duration(period))
		}
	}()
	go func() {
		defer wg.Done()
		for {
			logText := <-logCh
			println(logText)
		}
	}()
	wg.Wait()
}
