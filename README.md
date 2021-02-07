# Work Tracker
Track working time of your employees.

### Notes
* Only Linux based
* Need root privilages

### Features
* Time tracking
* Screenshots
* Keyboard keypress count

## Build
You need a Google Cloud service account for syncing data.
Without service account you can only use this offline.
- Create a Google Cloud service account.
- Create keys for this account and download as json.
- Set `DriveCredentials` in `constants/constants.go` to the contents of you json file.
- Create a folder in your Google Drive.
- Share this folder with write access to your service account using service account's email (something like `abc@xyz-123.iam.gserviceaccount.com`).
- Set `DriveBaseFolderId` in `constants/constants.go` to the id of the folder you just created.

To build the program, you can use:
```
go build
```
You will need root privileges to run it.
These are required to record keystrokes.
