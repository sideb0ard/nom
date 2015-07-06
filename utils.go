package main

import "time"

func percy(m int) int {
	return (m / 1000) * 100 // return m bytes as percentage of (hardcoded for the mo) eth0 speed 1GB
}

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
		time.Sleep(SLEEPTIME)
	}
}
