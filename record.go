package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type Record struct {
	WeeklyRecord    int
	DailyRecord     int
	Timestamp       time.Time
	ActiveWindow    string
	KeyboardStrokes int
	MouseStrokes    int
	ActivityLevel   int
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
