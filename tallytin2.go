package main
// Tally input, vmix output on raspberrypi
// Teppo Rekola 2019

import (
	"flag"
	"fmt"
  "time"

  "net/http"

	"github.com/kidoman/embd"

	_ "github.com/kidoman/embd/host/all"
)

func main() {
	flag.Parse()

	if err := embd.InitGPIO(); err != nil {
		panic(err)
	}
	defer embd.CloseGPIO()

  tally1, err := embd.NewDigitalPin(4)
	if err != nil {
		panic(err)
	}
	defer tally1.Close()

	if err := tally1.SetDirection(embd.In); err != nil {
		panic(err)
	}
	tally1.ActiveLow(false)

	err = tally1.Watch(embd.EdgeBoth, handleTally1)
	if err != nil {
		panic(err)
	}

	fmt.Printf("start")
  for true {
    time.Sleep(10000 * time.Millisecond)
    fmt.Println("ping")
  }

}
  func handleTally1(pin embd.DigitalPin) {
    pinValue, _ := pin.Read()
    if pinValue == 0 {
      fmt.Println("cut in")
      resp, err := http.Get("http://10.39.1.85:8088/api/?Function=MultiViewOverlayOn&Value=1&Input=6")
      if err != nil {
    		panic(err)
    	}
      resp.Body.Close()

    } else {
      fmt.Println("cut out")
      resp, err := http.Get("http://10.39.1.85:8088/api/?Function=MultiViewOverlayOff&Value=1&Input=6")
      if err != nil {
    		panic(err)
    	}
      resp.Body.Close()
    }

    time.Sleep(100 * time.Millisecond)
  }
