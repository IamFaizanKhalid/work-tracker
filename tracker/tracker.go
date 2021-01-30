package tracker

import (
	"github.com/kbinani/screenshot"
	"image"
	"log"
	"os/exec"
)

func GetActiveWindowName() string {
	b, err := exec.Command("xdotool", "getactivewindow", "getwindowname").Output()
	if err != nil {
		log.Printf("Error getting active window: %v", err)
		return "Undefined"
	}
	return string(b[:len(b)-1])
}

func GetScreenShot() *image.RGBA {
	img, err := screenshot.CaptureDisplay(0)
	if err != nil {
		log.Println("Error capturing screenshot: %v", err)
	}
	return img
}
