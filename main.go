package main

import (
	"fmt"
	"machine"
	"sync"
	"time"

	"git.o0.tel/sidc/unoblink/devices"
	"git.o0.tel/sidc/unoblink/types"
	"git.o0.tel/sidc/unoblink/view"
)

var state types.StateStage

func emitters() (types.TrinityLEDs, error) {
	led1, err := devices.NewLED(machine.D1, machine.TCC0)
	led2, err := devices.NewLED(machine.D2, machine.TCC1)
	led3, err := devices.NewLED(machine.D3, machine.TCC1)
	return types.TrinityLEDs{led1, led2, led3}, err
}

func termListener(wg *sync.WaitGroup, terminateChan <-chan bool, terminate *bool) {
	defer wg.Done()
	println("Started termination listener")
	die := <-terminateChan
	if die {
		*terminate = true
	}
}

func printStats(wg *sync.WaitGroup, printTick *time.Ticker, terminate *bool, leds types.TrinityLEDs) {
	defer wg.Done()
	for !*terminate {
		<-printTick.C
		fmt.Println("\033[H\033[2J")
		view.PrintStats(leds)
	}
}

func oldFn() {
	// for !terminate {
	// 	var internalWG sync.WaitGroup
	// 	for _, led := range emitters {
	// 		<-parallelRunDelay.C
	// 		internalWG.Add(1)
	// 		go led.Blink(&internalWG, pinChange)
	// 	}
	// 	internalWG.Wait()
	// 	runtime.GC()
	// }
}

type PinVolt struct {
	Pin     uint8
	Voltage uint32
}

func runner(rcChan chan *PinVolt, leds *types.TrinityLEDs, terminate *bool) {
	var pinWG sync.WaitGroup

	vRegulator := func(wg *sync.WaitGroup, ledNumber uint8, v chan<- *PinVolt) {
		defer wg.Done()
		defer fmt.Println("vRegulator ends!")
		pinChange := time.NewTicker(30 * time.Millisecond)
		defer pinChange.Stop()
		change := func(ledNumber uint8, volts uint32) {
			<-pinChange.C
			v <- &PinVolt{ledNumber, volts}
		}
		for volts := uint32(0); volts < 65000; volts += 650 {
			change(ledNumber, volts)
		}
		for volts := uint32(65000); volts > 0; volts -= 650 {
			change(ledNumber, volts)
		}
		change(ledNumber, 0)
		// runtime.Gosched()
	}
	vReceiver := func(rcChan <-chan *PinVolt, terminate *bool, leds *types.TrinityLEDs) {
		defer fmt.Println("vReceiver ends!")
		signalsReceived := 0
		for volts := range rcChan {
			if *terminate {
				break
			}
			signalsReceived++
			leds[volts.Pin].Set(volts.Voltage)
			fmt.Printf("writing led%d v: %d\n", volts.Pin, volts.Voltage)
		}
	}

	fmt.Println("starting vReceiver...")
	pinWG.Add(1)
	go vReceiver(rcChan, terminate, leds)

	pinWG.Add(1)
	go vRegulator(&pinWG, 0, rcChan)
	fmt.Println("vRegulator1 ready!")
	time.Sleep(1 * time.Second)

	pinWG.Add(1)
	go vRegulator(&pinWG, 1, rcChan)
	fmt.Println("vRegulator2 ready!")
	time.Sleep(1 * time.Second)

	pinWG.Add(1)
	go vRegulator(&pinWG, 2, rcChan)
	fmt.Println("vRegulator3 ready!")

	pinWG.Wait()
	close(rcChan)
}

func main() {
	fmt.Println("Starting...")
	time.Sleep(5 * time.Second)
	var mainWG sync.WaitGroup
	terminateChan := make(chan bool)
	terminate := false
	// printTick := time.NewTicker(250 * time.Millisecond)
	// defer printTick.Stop()

	emitters, err := emitters()
	if err != nil {
		panic(err)
	}
	rcChan := make(chan *PinVolt, 3)

	mainWG.Add(1)
	go termListener(&mainWG, terminateChan, &terminate)

	runner(rcChan, &emitters, &terminate)
	terminate = true

	fmt.Println("waiting main wg...")
	mainWG.Wait()
	fmt.Println("waited main wg: OK!")
}
