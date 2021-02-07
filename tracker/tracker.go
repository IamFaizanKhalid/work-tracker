package tracker

import (
	"github.com/IamFaizanKhalid/work-tracker/log"
	"github.com/MarinX/keylogger"
	"github.com/kbinani/screenshot"
	"image"
	"os/exec"
)

func GetActiveWindowName() string {
	b, err := exec.Command("xdotool", "getactivewindow", "getwindowname").Output()
	if err != nil {
		log.Error.Printf("Error getting active window: %v", err)
		return "Undefined"
	}
	return string(b[:len(b)-1])
}

func GetScreenShot() *image.RGBA {
	img, err := screenshot.CaptureDisplay(0)
	if err != nil {
		log.Error.Printf("Error capturing screenshot: %v\n", err)
	}
	return img
}

func GetKeyLogger() *keylogger.KeyLogger {
	kl, err := keylogger.New(keylogger.FindKeyboardDevice())
	if err != nil {
		log.Error.Printf("Error getting keylogger: %v\n", err)
		return nil
	}
	return kl
}
