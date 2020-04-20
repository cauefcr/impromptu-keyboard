package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	hook "github.com/cauefcr/ghook"
	"github.com/xlab/portmidi"
)

type MIDINote struct {
	Note     int64 `json:"note"`
	Velocity int64 `json:"velocity"`
	Duration int64 `json:"duration"`
}

type MIDIEvent struct {
	Status int64
	Note   MIDINote
}

func main() {
	EvChan := hook.Start()
	defer hook.End()
	m := map[string]MIDINote{}
	// low
	v, _ := strconv.Atoi(os.Args[1])
	off, _ := strconv.Atoi(os.Args[2])
	offset := int64(off) * 12
	vel := int64(v)
	m["q"] = MIDINote{Note: 0, Velocity: vel, Duration: 100}
	m["w"] = MIDINote{Note: 1, Velocity: vel, Duration: 100}
	m["e"] = MIDINote{Note: 2, Velocity: vel, Duration: 100}
	m["r"] = MIDINote{Note: 3, Velocity: vel, Duration: 100}
	m["t"] = MIDINote{Note: 4, Velocity: vel, Duration: 100}
	m["y"] = MIDINote{Note: 5, Velocity: vel, Duration: 100}
	m["u"] = MIDINote{Note: 6, Velocity: vel, Duration: 100}
	m["i"] = MIDINote{Note: 7, Velocity: vel, Duration: 100}
	m["o"] = MIDINote{Note: 8, Velocity: vel, Duration: 100}
	m["p"] = MIDINote{Note: 9, Velocity: vel, Duration: 100}
	// mid
	m["a"] = MIDINote{Note: 10, Velocity: vel, Duration: 100}
	m["s"] = MIDINote{Note: 11, Velocity: vel, Duration: 100}
	m["d"] = MIDINote{Note: 12, Velocity: vel, Duration: 100}
	m["f"] = MIDINote{Note: 13, Velocity: vel, Duration: 100}
	m["g"] = MIDINote{Note: 14, Velocity: vel, Duration: 100}
	m["h"] = MIDINote{Note: 15, Velocity: vel, Duration: 100}
	m["j"] = MIDINote{Note: 16, Velocity: vel, Duration: 100}
	m["k"] = MIDINote{Note: 17, Velocity: vel, Duration: 100}
	m["l"] = MIDINote{Note: 18, Velocity: vel, Duration: 100}
	//high
	m["z"] = MIDINote{Note: 18, Velocity: vel, Duration: 100}
	m["x"] = MIDINote{Note: 19, Velocity: vel, Duration: 100}
	m["c"] = MIDINote{Note: 20, Velocity: vel, Duration: 100}
	m["v"] = MIDINote{Note: 21, Velocity: vel, Duration: 100}
	m["b"] = MIDINote{Note: 22, Velocity: vel, Duration: 100}
	m["n"] = MIDINote{Note: 23, Velocity: vel, Duration: 100}
	m["m"] = MIDINote{Note: 24, Velocity: vel, Duration: 100}
	// pressed := map[string]bool{}

	fmt.Println("Starting sequence. Press Ctrl+C to quit...")

	// portmidi.Initialize()
	// defer portmidi.Terminate()
	portmidi.Initialize()
	defer portmidi.Terminate()
	id, _ := portmidi.DefaultOutputDeviceID()
	out, err := portmidi.NewOutputStream(id, 1024, 0, 1, portmidi.FilterUndefined)
	if err != nil {
		panic(err.Error())
	}

	sink := out.Sink()

	// fmt.Printf("%+v", portmidi.CountDevices())
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	on := true

	inst := 0
	// up instrument
	hook.Register(hook.KeyDown, []string{"ctrl", "9"}, func(e hook.Event) {
		fmt.Println(inst)
		sink <- portmidi.Event{Timestamp: int32(time.Now().Unix()), Message: portmidi.NewMessage(0xc0, byte(inst), 0)}
		inst++
		inst %= 127
	})
	// down instrument
	hook.Register(hook.KeyDown, []string{"ctrl", "8"}, func(e hook.Event) {
		fmt.Println(inst)
		sink <- portmidi.Event{Timestamp: int32(time.Now().Unix()), Message: portmidi.NewMessage(0xc0, byte(inst), 0)}
		inst--
		inst += 127
		inst %= 127
	})
	// up offset
	hook.Register(hook.KeyDown, []string{"ctrl", "7"}, func(e hook.Event) {
		offset += 12
	})
	// down offset
	hook.Register(hook.KeyDown, []string{"ctrl", "6"}, func(e hook.Event) {
		offset -= 12
	})
	// stop playing all notes
	hook.Register(hook.KeyDown, []string{"esc"}, func(e hook.Event) {
		for v := 0; v < 127; v++ {
			sink <- portmidi.Event{Timestamp: int32(time.Now().Unix()), Message: portmidi.NewMessage(0x80, byte(v), 0)}
		}
	})
	// toggle mute
	hook.Register(hook.KeyDown, []string{"ctrl", "0"}, func(e hook.Event) {
		on = !on
	})
	// quit
	hook.Register(hook.KeyDown, []string{"ctrl", "1"}, func(e hook.Event) {
		os.Exit(0)
	})
	for k, v := range m {
		val := v
		// key := k
		hook.Register(hook.KeyDown, []string{k}, func(e hook.Event) {
			if on {
				// !pressed[key]aqwaqweqqsdsaaaaqqasjkjkkkkkjhjkkjjaaasqwesssasdasasasasasasasasadadaqqdqdqdqdqdqdqdqdqdqdqdqdqqweqwqwqqweryuuiopasdfg
				fmt.Printf("\r%v    ", toNote(val.Note))
				sink <- portmidi.Event{Timestamp: int32(time.Now().Unix()), Message: portmidi.NewMessage(0x90, byte(val.Note+offset), byte(val.Velocity))}
				// pressed[key] = true
			}
		})
		hook.Register(hook.KeyUp, []string{k}, func(e hook.Event) {
			if on {
				// pressed[key]
				fmt.Println("plom")
				sink <- portmidi.Event{Timestamp: int32(time.Now().Unix()), Message: portmidi.NewMessage(0x80, byte(val.Note+offset), byte(val.Velocity))}
				// pressed[key] = false
			}
		})
	}

	go func() {
		<-hook.Process(EvChan)
	}()
	defer hook.End()
LOOP:
	for {
		select {
		case <-sigint:
			break LOOP
		}
	}
	out.Close()
}

func toNote(note int64) string {
	switch note % 12 {
	case 0:
		return "C"
	case 1:
		return "C#"
	case 2:
		return "D"
	case 3:
		return "D#"
	case 4:
		return "E"
	case 5:
		return "F"
	case 6:
		return "F#"
	case 7:
		return "G"
	case 8:
		return "G#"
	case 9:
		return "A"
	case 10:
		return "A#"
	case 11:
		return "B"
	}
	return ""
}
