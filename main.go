package main

import (
	"machine"
	"sync"
	"time"

	"git.o0.tel/sidc/unoblink/blinker"
)

// diginal pins { 3, 5, 6, 9, 10, 11 };
// Pins that works with PWM LEDs:
// 5,6

//	func listenToSignals(cancel context.CancelFunc) {
//		c := make(chan int, 1)
//		<-c
//		cancel()
//	}
//
//	func waitForDurationShutdown(cancel context.CancelFunc, dur time.Duration) {
//		<-time.After(dur)
//		cancel()
//	}

func main() {
	var wg sync.WaitGroup
	var pwm machine.PWM

	err := pwm.Configure(machine.PWMConfig{})
	if err != nil {
		println(err.Error())
		return
	}
	// ctx, _ := context.WithCancel(context.Background())
	println("creating blinkers")
	blinkers := []blinker.Blinker{
		// blinker.NewGracefulBlinker(pwm, machine.PD5),
		blinker.NewGracefulBlinker(blinker.NewLockedPWM(pwm), machine.PD6),
	}
	println("blinkers created")
	starter := blinker.NewFlawlessStarter()
	// starter.Ctx(ctx)
	starter.Blinkers(blinkers)
	starter.Delay(300 * time.Millisecond)
	println("starter configured")
	starter.Go(&wg)

	wg.Wait()
	println("waited")
}
