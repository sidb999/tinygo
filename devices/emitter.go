package devices

import (
	"machine"
	"sync"
	"time"
)

const (
	DefaultADCReference  = 3300
	DefaultADCSampleTime = 40 // use the longest acquisition time
	DefaultADCSamples    = 4
)

type LightEmitter struct {
	Pin machine.Pin
	ADC machine.ADC
	TCC *machine.TCC
	Ch  uint8
}

func (l *LightEmitter) Set(emitLvl int) {
	l.TCC.Set(l.Ch, uint32(emitLvl))
}

func (l *LightEmitter) Blink(wg *sync.WaitGroup, pinChange *time.Ticker) {
	defer wg.Done()
	for i := 0; i < 65000; i += 650 {
		<-pinChange.C
		l.Set(i)
	}
	for i := 65000; i > 0; i -= 650 {
		<-pinChange.C
		l.Set(i)
	}
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
