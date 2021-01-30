package main

import "time"

type Record struct {
	Timestamp       time.Time
	ActiveWindow    string
	KeyboardStrokes int
	MouseStrokes    int
	ActivityLevel   int
}
