package main

import (
	"math/rand"
	"strconv"
	"time"
)

func main() {
	for {
		if time.Now().Second() == 0 {
			// Each minute has one settlement tick.
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func getResult() string {
	n := rand.Intn(99999)
	s := strconv.Itoa(n)
	sLen := len(s)
	diffN := 5 - sLen
	for i := 0; i < diffN; i++ {
		s = "0" + s
	}
	return s
}
