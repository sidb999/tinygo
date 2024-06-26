package devices

import (
	"context"
	"fmt"
	"machine"
	"os"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/gammazero/deque"
)

type VoltageCalculator struct {
	volts    map[uint8]*deque.Deque[float32]
	ctx      context.Context
	pins     map[uint8]machine.ADC
	pollFreq int
}

// pollFreq is the frequency at which the voltage is polled in Hz.
func NewVoltageCalculator(ctx context.Context, pollFreq int, pins map[uint8]machine.ADC) *VoltageCalculator {
	vc := &VoltageCalculator{
		ctx:      ctx,
		pollFreq: pollFreq,
		volts:    make(map[uint8]*deque.Deque[float32], len(pins)),
		pins:     make(map[uint8]machine.ADC, len(pins)),
	}
	for pin, adc := range pins {
		vc.pins[pin] = adc
		vc.volts[pin] = deque.New[float32](0, pollFreq)
	}
	return vc
}

func (vc *VoltageCalculator) collectData(wg *sync.WaitGroup) {
	defer wg.Done()
	var min, max float32
	var minInt, maxInt uint16
	var lastADC uint16
	intervalDuration := time.Duration(int64(1000/vc.pollFreq)) * time.Millisecond
	pollInterval := time.NewTicker(intervalDuration)
	defer pollInterval.Stop()
	printTick := time.NewTicker(100 * time.Millisecond)
	defer printTick.Stop()
	min, max = 5., .0
	minInt, maxInt = 65534, 0
	refV := float32((maxVoltage * (DefaultADCReference / 1000)) / 65535)
	w := tabwriter.NewWriter(os.Stdout, 10, 1, 1, '\t', 0)
	for {
		select {
		case <-vc.ctx.Done():
			return
		case <-pollInterval.C:
			for pin, adc := range vc.pins {
				lastADC = adc.Get()
				if minInt > lastADC {
					minInt = lastADC
				}
				if maxInt < lastADC {
					maxInt = lastADC
				}
				voltage := float32(lastADC) * refV / float32(maxVoltage)
				if min > voltage {
					min = voltage
				}
				if max < voltage {
					max = voltage
				}
				if lastADC <= voltageStep*2 {
					voltage = .0
				}
				if vc.volts[pin].Len() >= vc.pollFreq {
					vc.volts[pin].PopFront()
				}
				vc.volts[pin].PushBack(voltage)
			}
		case <-printTick.C:
			fmt.Println("\033[H\033[2J")
			for pin, q := range vc.volts {
				var vAVG float32
				vAVG = .0
				for i := 0; i < q.Len(); i++ {
					vAVG += q.At(i)
				}
				vAVG = vAVG / float32(q.Len())
				fmt.Printf("D%d:\tCURR: %.2f\tAVG: %.2fV\n", pin, vc.volts[pin].PopFront(), vAVG)
			}
			fmt.Printf("min: %.2fV\tmax: %.2fV\n", min, max)
			fmt.Printf("min: %d\tmax: %d\n", minInt, maxInt)
			w.Flush()
		}
	}
}

func (vc *VoltageCalculator) Measure(wg *sync.WaitGroup) {
	wg.Add(1)
	go vc.collectData(wg)
}
