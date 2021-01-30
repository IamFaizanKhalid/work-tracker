package main

import (
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"time"
)

const DURATION = 10

func main() {
	capture()
}

func capture() {
	minutesPassed := 0
	captures := 0
	captureAfter := 1 + rand.Int()%DURATION
	ticker := time.NewTicker(time.Duration(captureAfter) * time.Second)

	for range ticker.C {
		captures += 1
		minutesPassed += captureAfter
		captureAfter = (1 + rand.Int()%DURATION) + (DURATION - captureAfter)
		ticker.Reset(time.Duration(captureAfter) * time.Second)

		x := Record{
			Timestamp:       time.Now(),
			ActiveWindow:    getActiveWindowName(),
			KeyboardStrokes: 0,
			MouseStrokes:    0,
			ActivityLevel:   0,
		}

		fmt.Printf("Logged today:\t%v\t\t\tLogged this week:\t%[1]v\n", getTimeLogged(captures))
		fmt.Printf("Active window:\t%s\n", x.ActiveWindow)
		fmt.Println()
	}
}

func getTimeLogged(captures int) string {
	t := captures * DURATION
	return fmt.Sprintf("%02d:%02d", t/60, t%60)
}

func getActiveWindowName() string {
	b, err := exec.Command("xdotool", "getactivewindow", "getwindowname").Output()
	if err != nil {
		log.Printf("Error getting active window: %v", err)
		return "Undefined"
	}
	return string(b[:len(b)-1])
}
