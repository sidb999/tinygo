package main

import (
	"machine"
	"sync"

	"git.o0.tel/sidc/unoblink/blinker"
)

var wg sync.WaitGroup

func main() {
	println("tickers")
	tickers := blinker.NewTickers()
	println("pwm")
	pwm1 := machine.TCC0
	err := pwm1.Configure(machine.PWMConfig{})
	if err != nil {
		println(err.Error())
	}
	pwm2 := machine.TCC1
	err = pwm2.Configure(machine.PWMConfig{})
	if err != nil {
		println(err.Error())
	}
	// ctx, _ := context.WithCancel(context.Background())
	blinkers := []*blinker.GracefulBlinker{
		blinker.NewGracefulBlinker(pwm1, machine.D1, tickers),
		blinker.NewGracefulBlinker(pwm2, machine.D2, tickers),
		// blinker.NewGracefulBlinker(pwm2, machine.D3, tickers),
	}
	// println("blinkers created")
	starter := blinker.NewFlawlessStarter(tickers)
	// starter.Ctx(ctx)
	starter.Blinkers(blinkers)
	println("starter configured")
	for {
		println("GO!")
		starter.Go(&wg)
		println("waiting")
		wg.Wait()
		println("waited")
	}
}
