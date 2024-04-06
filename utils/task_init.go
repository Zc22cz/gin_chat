package utils

import "time"

type TimeFunc func(interface{}) bool

func Timer(delay, tick time.Duration, fun TimeFunc, param interface{}) {
	go func() {
		if fun == nil {
			return
		}
		t := time.NewTimer(delay)
		for {
			select {
			case <-t.C:
				if fun(param) == false {
					return
				}
				t.Reset(tick)
			}
		}
	}()
}
