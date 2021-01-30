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

	minutesPassed := 0
	captureAfter := 1 + rand.Int()%DURATION
	ticker := time.NewTicker(time.Duration(captureAfter) * time.Second)
	defer ticker.Stop()

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

func getLastRecord() Record {
	today := time.Now()
	day := (6 + today.Weekday()) % 7 // 0: Monday, 6: Sunday

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
			return record
		}
		today.AddDate(0, 0, -1)
	}

	return Record{}
}
