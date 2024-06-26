package main

import (
	"context"
	"fmt"
	"machine"
	"sync"
	"time"

	"git.o0.tel/sidc/tinygo/devices"
	"git.o0.tel/sidc/tinygo/types"
)

func emitters() (devices.LEDArray, error) {
	led1, err := devices.NewLED(machine.D1, machine.TCC0)
	led2, err := devices.NewLED(machine.D2, machine.TCC1)
	led3, err := devices.NewLED(machine.D3, machine.TCC1)
	led4, err := devices.NewLED(machine.D5, machine.TCC0)
	return devices.LEDArray{led1, led2, led3, led4}, err
}

func main() {
	time.Sleep(5 * time.Second)
	fmt.Println("calc created")
	var mainWG sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	rcChan := make(chan *types.PinIntensity, 4)

	emitters, err := emitters()
	if err != nil {
		panic(err)
	}
	vc := devices.NewVoltageCalculator(ctx, 20, emitters.GetADCs())
	vc.Measure(&mainWG)

	emitters.Blink(rcChan, ctx)
	cancel()

	fmt.Println("waiting main wg...")
	mainWG.Wait()
	fmt.Println("waited main wg: OK!")
}
