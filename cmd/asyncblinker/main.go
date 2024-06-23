package main

import (
	"fmt"
	"machine"
	"runtime"
	"sync"
	"time"

	"git.o0.tel/sidc/tinygo/devices"
	"git.o0.tel/sidc/tinygo/types"
	"git.o0.tel/sidc/tinygo/view"
)

var state types.StateStage

func termListener(wg *sync.WaitGroup, terminateChan <-chan bool, terminate *bool) {
	defer wg.Done()
	println("Started termination listener")
	die := <-terminateChan
	if die {
		*terminate = true
	}
}

func printStats(wg *sync.WaitGroup, terminate *bool, leds devices.LEDArray) {
	printTick := time.NewTicker(50 * time.Millisecond)
	serialPrinter := view.NewSerialPrinter()
	defer printTick.Stop()
	defer wg.Done()
	for !*terminate {
		<-printTick.C
		fmt.Println("\033[H\033[2J")
		serialPrinter.PrintStats(leds)
	}
}

func emitters() (devices.LEDArray, error) {
	led1, err := devices.NewLED(machine.D1, machine.TCC0)
	led2, err := devices.NewLED(machine.D2, machine.TCC1)
	led3, err := devices.NewLED(machine.D3, machine.TCC1)
	led4, err := devices.NewLED(machine.D5, machine.TCC0)
	return devices.LEDArray{led1, led2, led3, led4}, err
}

func main() {
	runtime.GOMAXPROCS(10)
	var mainWG sync.WaitGroup
	terminateChan := make(chan bool)
	rcChan := make(chan *types.PinVolt, 3)
	terminate := false

	emitters, err := emitters()
	if err != nil {
		panic(err)
	}

	mainWG.Add(1)
	go termListener(&mainWG, terminateChan, &terminate)

	mainWG.Add(1)
	go printStats(&mainWG, &terminate, emitters)

	emitters.Blink(rcChan, &terminate)
	terminateChan <- true

	fmt.Println("waiting main wg...")
	mainWG.Wait()
	fmt.Println("waited main wg: OK!")
}
