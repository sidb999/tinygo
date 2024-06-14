package main

import (
	"fmt"
	"machine"
	"sync"
	"time"
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

//	func main() {
//		var wg sync.WaitGroup
//		pwm := blinker.NewLockedPWM()
//		// ctx, _ := context.WithCancel(context.Background())
//		println("creating blinkers")
//		blinkers := []blinker.Blinker{
//			blinker.NewGracefulBlinker(pwm, machine.D1),
//			blinker.NewGracefulBlinker(pwm, machine.D3),
//			blinker.NewGracefulBlinker(pwm, machine.D4),
//		}
//		println("blinkers created")
//		starter := blinker.NewFlawlessStarter()
//		// starter.Ctx(ctx)
//		starter.Blinkers(blinkers)
//		// starter.Delay(300 * time.Millisecond)
//		println("starter configured")
//		starter.Go(&wg)
//
//		wg.Wait()
//		println("waited")
//	}
type StateStage int8

var state StateStage

func (s StateStage) String() string {
	return [...]string{"Low pins", "High pins"}[s]
}

const (
	lowP StateStage = iota
	highP
)
const adcName = "LED"

func l(n int) string {
	return adcName + string(n)
}

func getVoltage(led machine.ADC) string {
	// vv := float32(led.Get()) * 3.3 * 1.1 * 2 / 65536 // 1.1 is just a kluge to get closer to expected battery at full charge
	// return fmt.Sprintf("%.2f (%x)", vv, v)
	return fmt.Sprintf("%.2f", float32(led.Get())*3.3/65536)
}

func printStats(stage StateStage, leds ...machine.ADC) {
	out := "\n" + stage.String()
	for n, led := range leds {
		out += "\n\t" + l(n) + ": " + getVoltage(led)
	}
	out += "\n"
	fmt.Println(out)
}

func main() {
	// printCh := make(chan bool)
	// switchCh := make(chan bool)
	println("init ADC")
	machine.InitADC() // init the machine's ADC subsystem
	pin1 := machine.D2
	pin1.Configure(machine.PinConfig{
		Mode: machine.PinOutput,
	})
	led1 := machine.ADC{Pin: pin1}

	pin2 := machine.D3
	pin2.Configure(machine.PinConfig{
		Mode: machine.PinOutput,
	})
	led2 := machine.ADC{Pin: pin2}

	pin3 := machine.D4
	pin3.Configure(machine.PinConfig{
		Mode: machine.PinOutput,
	})
	led3 := machine.ADC{Pin: pin3}

	var wg sync.WaitGroup
	wg.Add(2)
	printTick := time.NewTicker(500 * time.Millisecond)
	stateTick := time.NewTicker(5 * time.Second)
	pinChange := time.NewTicker(5 * time.Millisecond)

	go func() {
		defer wg.Done()
		for {
			// <-printCh
			<-printTick.C
			fmt.Println("\033[H\033[2J")
			printStats(state, led1, led2, led3)
		}
	}()
	go func() {
		defer wg.Done()
		for {
			led1.Pin.High()

			wg.Add(1)
			go func() {
				defer wg.Done()
				for i := 0; i < 400; i += 5 {
					<-pinChange.C
					led2.Configure(machine.ADCConfig{
						Reference:  uint32(i * 16),
						SampleTime: 2,
						Samples:    2,
					})
					led2.Pin.High()
					led3.Configure(machine.ADCConfig{
						Reference:  uint32(i * 16),
						SampleTime: 2,
						Samples:    2,
					})
				}
			}()

			state = highP
			<-stateTick.C

			led1.Pin.Low()

			wg.Add(1)

			go func() {
				defer wg.Done()
				for i := 400; i > 0; i -= 5 {
					<-pinChange.C
					led2.Configure(machine.ADCConfig{
						Reference:  uint32(i * 16),
						SampleTime: 2,
						Samples:    2,
					})
					led2.Pin.High()
					led3.Configure(machine.ADCConfig{
						Reference:  uint32(i * 16),
						SampleTime: 2,
						Samples:    2,
					})
				}
			}()

			state = lowP
			<-stateTick.C
		}
	}()
	// for {
	// 	time.Sleep(200 * time.Millisecond)
	// 	printCh <- true
	//
	// }
	wg.Wait()
}
