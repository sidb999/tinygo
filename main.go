package main

import (
	"fmt"
	"machine"
	"runtime"
	"sync"
	"time"

	"git.o0.tel/sidc/unoblink/devices"
	"git.o0.tel/sidc/unoblink/types"
	"git.o0.tel/sidc/unoblink/view"
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

const (
	maxVoltage  = 25000
	voltageStep = 1250
	startDelay  = 250 * time.Millisecond
)

type PinVolt struct {
	Pin     uint8
	Voltage uint32
}

func emitters() (types.TrinityLEDs, error) {
	led1, err := devices.NewLED(machine.D1, machine.TCC0)
	led2, err := devices.NewLED(machine.D2, machine.TCC1)
	led3, err := devices.NewLED(machine.D3, machine.TCC1)
	led4, err := devices.NewLED(machine.D5, machine.TCC0)
	return types.TrinityLEDs{led1, led2, led3, led4}, err
}

func runner(rcChan chan *PinVolt, leds *types.TrinityLEDs, terminate *bool) {
	var regulatorWG sync.WaitGroup
	var receiverWG sync.WaitGroup

	vRegulator := func(wg *sync.WaitGroup, ledNumber uint8, v chan<- *PinVolt) {
		defer wg.Done()
		defer fmt.Printf("vRegulator%d ends!\n", ledNumber)
		pinChange := time.NewTicker(30 * time.Millisecond)
		defer pinChange.Stop()
		change := func(ledNumber uint8, volts uint32) {
			<-pinChange.C
			v <- &PinVolt{ledNumber, volts}
		}
		for i := 0; i < 2; i++ {
			// for {
			for volts := uint32(0); volts < maxVoltage; volts += voltageStep {
				change(ledNumber, volts)
			}
			for volts := uint32(maxVoltage); volts > voltageStep; volts -= voltageStep {
				change(ledNumber, volts)
			}
			for volts := uint32(0); volts < maxVoltage-uint32(maxVoltage/5); volts += voltageStep {
				change(ledNumber, 0)
			}
		}
		change(ledNumber, 0)
	}
	vReceiver := func(wg *sync.WaitGroup, rcChan <-chan *PinVolt, terminate *bool, leds *types.TrinityLEDs) {
		defer wg.Done()
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
	receiverWG.Add(1)
	go vReceiver(&receiverWG, rcChan, terminate, leds)
	time.Sleep(startDelay)

	for i := uint8(0); i < uint8(len(*leds)); i++ {
		regulatorWG.Add(1)
		go vRegulator(&regulatorWG, i, rcChan)
		time.Sleep(startDelay)
	}

	regulatorWG.Wait()
	fmt.Println("regulators waited: OK!")
	close(rcChan)
	receiverWG.Wait()
	fmt.Println("receivers waited: OK!")
}

func main() {
	runtime.GOMAXPROCS(10)
	fmt.Println("Starting...")
	var mainWG sync.WaitGroup
	terminateChan := make(chan bool)
	terminate := false

	emitters, err := emitters()
	if err != nil {
		panic(err)
	}
	rcChan := make(chan *PinVolt, 3)

	mainWG.Add(1)
	go termListener(&mainWG, terminateChan, &terminate)

	runner(rcChan, &emitters, &terminate)
	terminateChan <- true

	fmt.Println("waiting main wg...")
	mainWG.Wait()
	fmt.Println("waited main wg: OK!")
}
