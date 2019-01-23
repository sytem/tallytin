package main

// Tally input, vmix output on raspberrypi
// Teppo Rekola 2019

import (
	"flag"
	"fmt"
	"time"

	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/all"
	"net/http"
)

func main() {
	flag.Parse()

	if err := embd.InitGPIO(); err != nil {
		panic(err)
	}
	defer embd.CloseGPIO()

	// first input
	tally1, err := embd.NewDigitalPin(02)
	if err != nil {
		panic(err)
	}
	defer tally1.Close()

	if err := tally1.SetDirection(embd.In); err != nil {
		panic(err)
	}
	tally1.ActiveLow(false)
	tally1.PullUp()
	
	err = tally1.Watch(embd.EdgeBoth, handleTally1)
	if err != nil {
		panic(err)
	}

	//second input
	tally2, err := embd.NewDigitalPin(04)
	if err != nil {
		panic(err)
	}
	defer tally2.Close()

	if err := tally2.SetDirection(embd.In); err != nil {
		panic(err)
	}
	tally2.ActiveLow(false)
	tally2.PullUp()

	err = tally2.Watch(embd.EdgeBoth, handleTally2)
	if err != nil {
		panic(err)
	}

	fmt.Printf("start")
	for {
		time.Sleep(10000 * time.Millisecond)
		fmt.Println("ping")
	}

}

func handleTally1(pin embd.DigitalPin) {
	pinValue, _ := pin.Read()
	if pinValue == 0 {
		fmt.Println("1 cut in")
		resp, err := http.Get("http://10.39.1.85:8088/api/?Function=MultiViewOverlayOn&Value=1&Input=6")
		if err != nil {
			panic(err)
		}
		resp.Body.Close()

	} else {
		fmt.Println("1 cut out")
		resp, err := http.Get("http://10.39.1.85:8088/api/?Function=MultiViewOverlayOff&Value=1&Input=6")
		if err != nil {
			panic(err)
		}
		resp.Body.Close()
	}

	time.Sleep(100 * time.Millisecond)
}

func handleTally2(pin embd.DigitalPin) {
	pinValue, _ := pin.Read()
	if pinValue == 0 {
		fmt.Println("2 cut in")
		resp, err := http.Get("http://10.39.1.85:8088/api/?Function=Transition1")
		if err != nil {
			panic(err)
		}
		resp.Body.Close()

	} else {
		fmt.Println("2 cut out")
		resp, err := http.Get("http://10.39.1.85:8088/api/?Function=Transition2")
		if err != nil {
			panic(err)
		}
		resp.Body.Close()
	}

	time.Sleep(100 * time.Millisecond)
}
