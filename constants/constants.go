package constants

import "time"

var (
	WorkDir    string
	CurrentDir string
)

const (
	DURATION      = 10
	DURATION_UNIT = time.Minute
	MIN_ACTIVITY  = 50

	DriveBaseFolderId = "" // Your folder id from google drive
	DriveCredentials  = "" // Your service account credentials from json file

	WorkLogFileName   = "logs"
	WorkLogIdFileName = "logs.id"
	DriveIdFileName   = "drive.id"

	LogFileName = "logfile"
)
