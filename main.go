package main

import (
	"fmt"
	"github.com/IamFaizanKhalid/work-tracker/tracker"
	"image/png"
	"log"
	"math/rand"
	"os"
	"time"
)

const DURATION = 10

var WorkDir string
var CaptureDir string

type Record struct {
	Timestamp       time.Time
	ActiveWindow    string
	KeyboardStrokes int
	MouseStrokes    int
	ActivityLevel   int
}

func main() {
	HomeDirectory, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error getting user home directory: %v", err)
	}

	WorkDir = HomeDirectory + "/.work-tracker"
	CaptureDir = WorkDir + "/captures"

	err = os.MkdirAll(CaptureDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Error creating working directory: %v", err)
	}

	startTracking()
}

func startTracking() {
	minutesPassed := 0
	captures := 0
	captureAfter := 1 + rand.Int()%DURATION
	ticker := time.NewTicker(time.Duration(captureAfter) * time.Minute)

	for range ticker.C {
		captures += 1
		minutesPassed += captureAfter
		captureAfter = (1 + rand.Int()%DURATION) + (DURATION - captureAfter)
		ticker.Reset(time.Duration(captureAfter) * time.Minute)

		timestamp := time.Now()
		saveScreenshot(timestamp)

		x := Record{
			Timestamp:       timestamp,
			ActiveWindow:    tracker.GetActiveWindowName(),
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

func saveScreenshot(timestamp time.Time) {
	fileName := fmt.Sprintf(CaptureDir+"/capture_%s.png", timestamp.Format("20060102150405"))

	file, err := os.Create(fileName)
	if err != nil {
		log.Printf("Error creating file: %v", err)
	}
	defer file.Close()

	err = png.Encode(file, tracker.GetScreenShot())
	if err != nil {
		log.Printf("Error writing image: %v", err)
	}
}
