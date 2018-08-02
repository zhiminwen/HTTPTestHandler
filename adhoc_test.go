package main

import (
	"log"
	"testing"
	"time"
)

func TestFloat(t *testing.T) {
	var load float32
	load = 0.8
	log.Printf("ns=%d", time.Now().UnixNano())
	log.Printf("us=%d", time.Now().UnixNano()/1e+3)
	log.Printf("ms=%d", time.Now().UnixNano()/1e+6)
	log.Printf("s=%d", time.Now().Unix())
	log.Printf("duration = %v", time.Duration((1-load)*100)*time.Millisecond)
}
