package devices

import (
	"context"
	"errors"
	"fmt"
	"machine"
	"sync"
	"time"

	"git.o0.tel/sidc/tinygo/types"
)

const (
	DefaultADCReference  = 3050
	DefaultADCResolution = 12
	DefaultADCSamples    = 32
	DefaultADCSampleTime = 40 // use the longest acquisition time
	// maxVoltage           = 25000
	maxVoltage  = 65535
	voltageStep = 3275
	startDelay  = 250 * time.Millisecond
)

type LightEmitter struct {
	Pin machine.Pin
	ADC machine.ADC
	TCC *machine.TCC
	Ch  uint8
}

func (l LightEmitter) Set(emitLvl uint32) {
	l.TCC.Set(l.Ch, uint32(emitLvl))
}

func NewLED(p machine.Pin, tcc *machine.TCC) (*LightEmitter, error) {
	var err error
	emitter := &LightEmitter{}
	emitter.Pin = p
	emitter.Pin.Configure(machine.PinConfig{
		Mode: machine.PinOutput,
	})

	emitter.ADC = machine.ADC{Pin: emitter.Pin}
	emitter.ADC.Configure(machine.ADCConfig{
		Reference:  DefaultADCReference,
		Resolution: DefaultADCResolution,
		Samples:    DefaultADCSamples,
		SampleTime: DefaultADCSampleTime,
	})
	emitter.TCC = tcc
	err = emitter.TCC.Configure(machine.PWMConfig{})
	if err != nil {
		return emitter, err
	}
	emitter.Ch, err = emitter.TCC.Channel(emitter.Pin)
	if err != nil {
		return emitter, err
	}
	return emitter, err
}

type LEDArray [4]*LightEmitter

func (leds LEDArray) Blink(rcChan chan *types.PinIntensity, ctx context.Context) {
	var regulatorWG sync.WaitGroup
	var receiverWG sync.WaitGroup

	vRegulator := func(wg *sync.WaitGroup, ledNumber uint8, v chan<- *types.PinIntensity) {
		defer wg.Done()
		defer fmt.Printf("\nvRegulator%d ends!\n", ledNumber)
		pinChange := time.NewTicker(30 * time.Millisecond)
		defer pinChange.Stop()
		change := func(ledNumber uint8, volts uint32, ctx context.Context) error {
			<-pinChange.C
			select {
			case <-ctx.Done():
				var err error
				err = ctx.Err()
				if err == context.DeadlineExceeded {
					err = errors.New("deadline exceeded!")
				}
				if err == context.Canceled {
					err = errors.New("canceled!")
				}
				fmt.Println(err.Error())
				return err
			default:
				v <- &types.PinIntensity{Pin: ledNumber, Voltage: volts}
				return nil
			}
		}
		// for i := 0; i < 10; i++ {
		for {
			for volts := uint32(0); volts < maxVoltage; volts += voltageStep {
				err := change(ledNumber, volts, ctx)
				if err != nil {
					return
				}
			}
			for volts := uint32(maxVoltage); volts > voltageStep; volts -= voltageStep {
				err := change(ledNumber, volts, ctx)
				if err != nil {
					return
				}
			}
			for volts := uint32(0); volts < maxVoltage-uint32(maxVoltage/5); volts += voltageStep {
				err := change(ledNumber, 0, ctx)
				if err != nil {
					return
				}
			}
		}
	}
	vReceiver := func(wg *sync.WaitGroup, rcChan <-chan *types.PinIntensity, ctx context.Context, leds *LEDArray) {
		defer wg.Done()
		defer fmt.Println("\nvReceiver ends!")
		for volts := range rcChan {
			select {
			case <-ctx.Done():
				return
			default:
				leds[volts.Pin].Set(volts.Voltage)
			}
		}
	}

	fmt.Println("starting vReceiver...")
	receiverWG.Add(1)
	go vReceiver(&receiverWG, rcChan, ctx, &leds)
	time.Sleep(startDelay)

	for i := uint8(0); i < uint8(len(leds)); i++ {
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

func (leds LEDArray) GetADCs() map[uint8]machine.ADC {
	adcMap := make(map[uint8]machine.ADC)
	for _, led := range leds {
		adcMap[uint8(led.Pin)] = led.ADC
	}
	return adcMap
}
