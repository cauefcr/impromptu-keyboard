package main

import (
	"fmt"

	hook "github.com/cauefcr/ghook"
	"github.com/progrium/go-shell"
)

func main() {
	EvChan := hook.Start()
	defer hook.End()
	var sh = shell.Run
	m := map[string]int{}
	// low
	m["q"] = 130
	m["w"] = 146
	m["e"] = 146
	m["r"] = 174
	m["t"] = 196
	m["y"] = 220
	m["u"] = 247
	m["i"] = 261
	m["o"] = 293
	m["p"] = 329
	// mid
	m["a"] = 261
	m["s"] = 294
	m["d"] = 329
	m["f"] = 349
	m["g"] = 392
	m["h"] = 440
	m["j"] = 493
	m["k"] = 523
	m["l"] = 587
	//high
	m["z"] = 523
	m["x"] = 587
	m["c"] = 659
	m["v"] = 698
	m["b"] = 783
	m["n"] = 880
	m["m"] = 988
	// m["<"] = 987
	// m[">"] = 1046
	// m[";"] = 1174
	for k, v := range m {
		val := v
		hook.Register(hook.KeyDown, []string{"alt", k}, func(e hook.Event) {
			go sh("play -n synth triangle " + fmt.Sprint(float32(val)*2) + " square " + fmt.Sprint(float32(val)*4) + " sin " + fmt.Sprint(float32(val)*8) + " gain -2.4 fade h 0.008 remix - fade h 0.1 0.5 1.5")
		})
	}

	<-hook.Process(EvChan)
}
