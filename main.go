package main

import (
	"fmt"
	"machine"
	"sync"
	"time"

	"git.o0.tel/sidc/unoblink/devices"
	"git.o0.tel/sidc/unoblink/view"
)

var state view.StateStage

func main() {
	var mainWG sync.WaitGroup
	errChan := make(chan error)
	terminateChan := make(chan bool)
	terminate := false
	pinChange := time.NewTicker(10 * time.Millisecond)
	printTick := time.NewTicker(250 * time.Millisecond)
	parallelRunDelay := time.NewTicker(250 * time.Millisecond)
	defer pinChange.Stop()
	defer printTick.Stop()
	defer parallelRunDelay.Stop()

	led1, err := devices.NewLED(machine.D1, machine.TCC0) // TCC0 channel 0
	if err != nil {
		// TODO: process real typed error and exit if critical
		errChan <- err
	}
	led2, err := devices.NewLED(machine.D2, machine.TCC1) // TCC1 channel 0
	if err != nil {
		errChan <- err
	}

	led3, err := devices.NewLED(machine.D3, machine.TCC1) // TCC1/0 channel 1/3
	if err != nil {
		errChan <- err
	}
	emitters := []*devices.LightEmitter{led1, led2, led3}

	mainWG.Add(1)
	go func() {
		defer mainWG.Done()
		println("Started termination listener")
		die := <-terminateChan
		if die {
			terminate = true
		}
	}()

	mainWG.Add(1)
	go func() {
		defer mainWG.Done()
		for !terminate {
			<-printTick.C
			fmt.Println("\033[H\033[2J")
			view.PrintStats(state, led1, led2, led3)
		}
	}()

	for !terminate {
		var internalWG sync.WaitGroup
		for _, led := range emitters {
			<-parallelRunDelay.C
			internalWG.Add(1)
			go led.Blink(&internalWG, pinChange)
		}
		internalWG.Wait()
	}
	mainWG.Wait()
}
