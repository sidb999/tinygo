package devices

import (
	"fmt"
	"machine"
	"sync"
	"time"

	"git.o0.tel/sidc/tinygo/types"
)

const (
	DefaultADCReference  = 3300
	DefaultADCSampleTime = 40 // use the longest acquisition time
	DefaultADCSamples    = 4
	maxVoltage           = 25000
	voltageStep          = 1250
	startDelay           = 250 * time.Millisecond
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
		SampleTime: DefaultADCSampleTime,
		Samples:    DefaultADCSamples,
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

func (leds LEDArray) Blink(rcChan chan *types.PinVolt, terminate *bool) {
	var regulatorWG sync.WaitGroup
	var receiverWG sync.WaitGroup

	vRegulator := func(wg *sync.WaitGroup, ledNumber uint8, v chan<- *types.PinVolt) {
		defer wg.Done()
		defer fmt.Printf("vRegulator%d ends!\n", ledNumber)
		pinChange := time.NewTicker(30 * time.Millisecond)
		defer pinChange.Stop()
		change := func(ledNumber uint8, volts uint32) {
			<-pinChange.C
			v <- &types.PinVolt{Pin: ledNumber, Voltage: volts}
		}
		for i := 0; i < 20; i++ {
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
	vReceiver := func(wg *sync.WaitGroup, rcChan <-chan *types.PinVolt, terminate *bool, leds *LEDArray) {
		defer wg.Done()
		defer fmt.Println("vReceiver ends!")
		for volts := range rcChan {
			if *terminate {
				break
			}
			leds[volts.Pin].Set(volts.Voltage)
		}
	}

	fmt.Println("starting vReceiver...")
	receiverWG.Add(1)
	go vReceiver(&receiverWG, rcChan, terminate, &leds)
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
