package view

import (
	"fmt"

	"git.o0.tel/sidc/tinygo/devices"
	"git.o0.tel/sidc/tinygo/types"
)

const (
	LowP types.StateStage = iota
	HighP
)
const ledStr = "LED"

type MachineADC interface {
	Get() uint16
}

type ledN string

func (l ledN) String(num uint8) string {
	return fmt.Sprintf("%s%d", l, num)
}

type voltageStat struct {
	stats [4]float32
}

func (v *voltageStat) Push(voltage float32) {
	v.stats[0] = v.stats[1]
	v.stats[1] = v.stats[2]
	v.stats[2] = v.stats[3]
	v.stats[3] = voltage
}

func NewVoltageStat() *voltageStat {
	return &voltageStat{stats: [4]float32{0, 0, 0, 0}}
}

type PinStats struct {
	pins map[uint8]*voltageStat
}

func (p *PinStats) Add(pin uint8, voltage float32) {
	if p.pins[pin] == nil {
		p.pins[pin] = NewVoltageStat()
		p.pins[pin].stats[3] = voltage
		return
	}
	p.pins[pin].Push(voltage)
}

func (p *PinStats) Get(pin uint8) float32 {
	vs := p.pins[pin]
	return float32((vs.stats[0] + vs.stats[1] + vs.stats[2] + vs.stats[3]) / 4)
}

type SerialPrinter struct {
	pinstats *PinStats
}

func NewSerialPrinter() *SerialPrinter {
	return &SerialPrinter{pinstats: &PinStats{pins: make(map[uint8]*voltageStat)}}
}

func (sp *SerialPrinter) GetVoltage(pin uint8, led MachineADC) string {
	sp.pinstats.Add(pin, float32(led.Get()))
	return fmt.Sprintf("%.2f", sp.pinstats.Get(pin)*3.3/65536)
}

func (sp *SerialPrinter) PrintStats(leds devices.LEDArray) {
	out := "\n"
	for n, led := range leds {
		out += fmt.Sprintf("\n\t%s: %s", ledN(ledStr).String(uint8(n)), sp.GetVoltage(uint8(n), led.ADC))
	}
	out += "\n"
	fmt.Println(out)
}
