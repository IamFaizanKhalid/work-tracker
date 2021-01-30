package main

import (
	"bytes"
	"encoding/json"
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
var CurrentDir string

type Record struct {
	WeeklyRecord    int
	DailyRecord     int
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
	CurrentDir = WorkDir + "/" + time.Now().Format("20060102")

	err = os.MkdirAll(CurrentDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Error creating working directory: %v", err)
	}

	startTracking()
}

func startTracking() {
	var record Record

	minutesPassed := 0
	captureAfter := 1 + rand.Int()%DURATION
	ticker := time.NewTicker(time.Duration(captureAfter) * time.Second)

	for range ticker.C {
		record.DailyRecord += 1
		record.WeeklyRecord += 1
		minutesPassed += captureAfter
		captureAfter = (1 + rand.Int()%DURATION) + (DURATION - captureAfter)
		ticker.Reset(time.Duration(captureAfter) * time.Second)

		record.Timestamp = time.Now()
		saveScreenshot(record.Timestamp)
		record.ActiveWindow = tracker.GetActiveWindowName()

		record.log()
		record.print()
	}
}

func getTimeLogged(captures int) string {
	t := captures * DURATION
	return fmt.Sprintf("%02d:%02d", t/60, t%60)
}

func saveScreenshot(timestamp time.Time) {
	fileName := fmt.Sprintf(CurrentDir+"/%s.png", timestamp.Format("150405"))

	file, err := os.Create(fileName)
	if err != nil {
		log.Printf("Error creating file: %v", err)
	}
	defer file.Close()

	encoder := &png.Encoder{
		CompressionLevel: png.BestCompression,
	}

	err = encoder.Encode(file, tracker.GetScreenShot())
	if err != nil {
		log.Printf("Error writing image: %v", err)
	}
}

func (r *Record) log() {
	b, err := json.Marshal(r)
	if err != nil {
		log.Printf("Error marshalling record: %v", err)
	}

	buffer := new(bytes.Buffer)
	err = json.Compact(buffer, b)
	if err != nil {
		log.Printf("Error compacting json record: %v", err)
	}
	buffer.WriteByte('\n')

	file, err := os.OpenFile(CurrentDir+"/logs", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0660)
	if err != nil {
		log.Printf("Error opening log file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(buffer.Bytes())
	if err != nil {
		log.Printf("Error writing log file: %v", err)
	}
}

func (r *Record) print() {
	fmt.Printf("%s\n--\n", r.Timestamp.Format("Monday, 02 Jan 2006 15:04:05 MST"))
	fmt.Printf("> Logged this week:\t%v\n", getTimeLogged(r.WeeklyRecord))
	fmt.Printf("> Logged today:\t\t%v\n", getTimeLogged(r.DailyRecord))
	fmt.Printf("> Active window:\t%s\n", r.ActiveWindow)
	fmt.Printf("> Activity level:\t%d\n", r.ActivityLevel)
	fmt.Printf("> Events:\t\t%d keyboard, %d mouse\n", r.KeyboardStrokes, r.MouseStrokes)
	fmt.Printf("------------------------------\n\n")
}
