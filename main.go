package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/IamFaizanKhalid/work-tracker/tracker"
	"image/png"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

const DURATION = 10

var WorkDir string
var CurrentDir string

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
	record := getLastRecord()

	// Ticker to trigger capture
	rand.Seed(time.Now().UTC().UnixNano())
	captureAfter := 1 + rand.Int()%DURATION
	ticker := time.NewTicker(time.Duration(captureAfter) * time.Minute)
	defer ticker.Stop()

	// Ticker to change day
	n := time.Now()
	d := time.Until(time.Date(n.Year(), n.Month(), n.Day()+1, 0, 0, 0, 0, n.Location()))
	dayChange := time.NewTicker(d)
	defer dayChange.Stop()

	// Channel to detect interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Get KeyLogger
	keyLogger := tracker.GetKeyLogger()
	if keyLogger == nil {
		return
	}

	fmt.Println("Time tracking started..")
	for {
		select {
		case now := <-ticker.C:
			record.DailyRecord += 1
			record.WeeklyRecord += 1
			captureAfter = (1 + rand.Int()%DURATION) + (DURATION - captureAfter)
			ticker.Reset(time.Duration(captureAfter) * time.Minute)

			record.Timestamp = now
			record.ActiveWindow = tracker.GetActiveWindowName()
			saveScreenshot(now)

			record.log()
			record.print()

			record.KeyboardStrokes = 0
			record.MouseStrokes = 0

		case now := <-dayChange.C:
			CurrentDir = WorkDir + "/" + now.Format("20060102")

			err := os.MkdirAll(CurrentDir, os.ModePerm)
			if err != nil {
				log.Fatalf("Error creating working directory: %v", err)
			}

			dayChange.Reset(24 * time.Hour)

			record.DailyRecord = 0
			if now.Weekday() == time.Monday {
				record.WeeklyRecord = 0
			}

		case e := <-keyLogger.Read():
			if e.KeyPress() {
				record.KeyboardStrokes++
			}

		case <-c:
			fmt.Println("\nTime tracking stopped..")
			return
		}
	}
}

func getTimeLogged(captures int) string {
	t := captures * DURATION
	return fmt.Sprintf("%02d:%02d", t/60, t%60)
}

func saveScreenshot(timestamp time.Time) {
	fileName := fmt.Sprintf(CurrentDir+"/%s", timestamp.Format("150405"))

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

func getLastRecord() Record {
	today := time.Now()
	currentDay := (6 + today.Weekday()) % 7 // 0: Monday, 6: Sunday
	day := currentDay

	for ; day >= 0; day-- {
		dir := WorkDir + "/" + today.Format("20060102")

		file, err := os.Open(dir + "/logs")
		if err == nil {
			scanner := bufio.NewScanner(file)
			var lastText string
			for scanner.Scan() {
				lastText = scanner.Text()
			}

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}

			if lastText == "" {
				return Record{}
			}

			var record Record
			err = json.Unmarshal([]byte(lastText), &record)
			if err != nil {
				log.Printf("Error getting last record: %v\n", err)
				return Record{}
			}

			record.KeyboardStrokes = 0
			record.MouseStrokes = 0
			if day != currentDay {
				record.DailyRecord = 0
			}

			return record
		}
		today = today.AddDate(0, 0, -1)
	}

	return Record{}
}
