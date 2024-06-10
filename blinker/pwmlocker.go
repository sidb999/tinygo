package blinker

import (
	"machine"
	"sync"
)

type LockedPWM struct {
	lock sync.Mutex
	pwm  machine.PWM
}

func (p *LockedPWM) Get() machine.PWM {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.pwm
}

func NewLockedPWM(pwm machine.PWM) *LockedPWM {
	return &LockedPWM{
		pwm: pwm,
	}
}
