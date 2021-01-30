package tracker

import (
	"fmt"
	"github.com/MarinX/keylogger"
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

func GetKeyLogger() *keylogger.KeyLogger {
	fmt.Println(keylogger.FindAllKeyboardDevices())
	kl, err := keylogger.New(keylogger.FindKeyboardDevice())
	if err != nil {
		log.Printf("Error getting keylogger: %v\n", err)
		return nil
	}
	return kl
}
