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
	Voltage int
}

func main() {
	fmt.Println("Starting...")
	time.Sleep(5 * time.Second)
	var mainWG sync.WaitGroup
	terminateChan := make(chan bool)
	terminate := false
	pinChange := time.NewTicker(30 * time.Millisecond)
	printTick := time.NewTicker(250 * time.Millisecond)
	parallelRunDelay := time.NewTicker(250 * time.Millisecond)
	defer pinChange.Stop()
	defer printTick.Stop()
	defer parallelRunDelay.Stop()

	emitters, err := emitters()
	if err != nil {
		panic(err)
	}
	var vChans types.TrinityChans
	mainChan := make(chan *PinVolt, 3)

	vRegulator := func() <-chan int {
		<-parallelRunDelay.C
		v := make(chan int, 1)
		go func() {
			defer close(v)
			for i := 0; i < 65000; i += 650 {
				v <- i
			}
			for i := 65000; i > 0; i -= 650 {
				v <- i
			}
			fmt.Println("vRegulator ends!")
		}()
		return v
	}
	vReceiver := func(volts *types.TrinityChans, terminate *bool, mainChan chan<- *PinVolt) {
		go func() {
			defer close(mainChan)
			defer fmt.Println("vReceiver ends!")
			signalsReceived := 0
			for !*terminate {
				select {
				case v1 := <-volts[0]:
					signalsReceived++
					<-pinChange.C
					fmt.Printf("received led%d v: %d\n", 0, v1)
					mainChan <- &PinVolt{0, v1}
				case v2 := <-volts[1]:
					signalsReceived++
					<-pinChange.C
					fmt.Printf("received led%d v: %d\n", 1, v2)
					mainChan <- &PinVolt{1, v2}
				case v3 := <-volts[2]:
					signalsReceived++
					<-pinChange.C
					fmt.Printf("received led%d v: %d\n", 2, v3)
					mainChan <- &PinVolt{2, v3}
				}
				// if signalsReceived > 10 {
				// 	fmt.Printf("received %d of signals\n", signalsReceived)
				// 	*terminate = true
				// }
			}
		}()
	}
	vWriter := func(mainChan <-chan *PinVolt, leds *types.TrinityLEDs, terminate *bool) {
		fmt.Println("vWriter starts!")
		for v := range mainChan {
			leds[v.Pin].Set(v.Voltage)
			fmt.Printf("writing led%d v: %d\n", v.Pin, v.Voltage)
		}
		fmt.Println("vWriter ends!")
	}

	fmt.Println("preloading handlers...")
	vChans[0] = vRegulator()
	fmt.Println("vRegulator1 ready!")
	vChans[1] = vRegulator()
	fmt.Println("vRegulator2 ready!")
	vChans[2] = vRegulator()
	fmt.Println("vRegulator3 ready!")
	fmt.Println("chans filled!")

	mainWG.Add(1)
	go termListener(&mainWG, terminateChan, &terminate)

	mainWG.Add(1)
	// go printStats(&mainWG, printTick, &terminate, emitters)
	fmt.Println("starting vReceiver")
	vReceiver(&vChans, &terminate, mainChan)
	fmt.Println("starting vWriter")
	vWriter(mainChan, &emitters, &terminate)

	fmt.Println("waiting main wg")
	mainWG.Wait()
}
