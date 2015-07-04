package main

import "time"

func reverseMap(m map[string]string) map[string]string {
	n := make(map[string]string)
	for k, v := range m {
		n[v] = k
	}
	return n
}

func oldTimer(timerChannel chan int) {
	timerChannel <- 1
	for {
		timerChannel <- 1
		time.Sleep(time.Second / 2)
	}
}
