package uploader

import (
	"github.com/IamFaizanKhalid/work-tracker/constants"
	"github.com/IamFaizanKhalid/work-tracker/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"io"
	"io/ioutil"
	"os"
	"time"
)

func getService() (*drive.Service, error) {
	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.JWTConfigFromJSON([]byte(constants.DriveCredentials), drive.DriveScope)
	if err != nil {
		log.Error.Println(err)
		return nil, err
	}

	client := config.Client(oauth2.NoContext)

	service, err := drive.New(client)
	if err != nil {
		log.Error.Printf("Cannot create the Google Drive service: %v\n", err)
		return nil, err
	}

	return service, err
}

func createDir(service *drive.Service, name string, parentId string) (*drive.File, error) {
	d := &drive.File{
		Name:     name,
		MimeType: "application/vnd.google-apps.folder",
		Parents:  []string{parentId},
	}

	file, err := service.Files.Create(d).Do()

	if err != nil {
		log.Error.Println(err)
		return nil, err
	}

	return file, nil
}

func uploadImage(service *drive.Service, name string, content io.Reader, parentId string) (*drive.File, error) {
	f := &drive.File{
		MimeType: "image/png",
		Name:     name,
		Parents:  []string{parentId},
	}
	file, err := service.Files.Create(f).Media(content).Do()

	if err != nil {
		log.Error.Println(err)
		return nil, err
	}

	return file, nil
}

func updateLogFile(service *drive.Service, localParentPath string, driveParentId string) (*drive.File, error) {
	f := &drive.File{
		MimeType: "text/plain",
		Name:     constants.WorkLogFileName,
	}

	var driveFile *drive.File
	logFile := localParentPath + "/" + constants.WorkLogFileName
	logIdFile := localParentPath + "/" + constants.WorkLogIdFileName

	content, err := os.Open(logFile)
	if err != nil {
		log.Error.Println(err)
		return nil, err
	}

	b, err := ioutil.ReadFile(logIdFile)
	if err != nil {
		f.Parents = []string{driveParentId}
		driveFile, err = service.Files.Create(f).Media(content).Do()
		if err != nil {
			log.Error.Println(err)
			return nil, err
		}
	} else {
		fileId := string(b)
		driveFile, err = service.Files.Update(fileId, f).Media(content).Do()
		if err != nil {
			log.Error.Println(err)
			return nil, err
		}
	}

	idFile, err := os.OpenFile(logIdFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Error.Println(err)
		return nil, err
	}

	_, err = idFile.WriteString(driveFile.Id)
	if err != nil {
		log.Error.Println(err)
		return nil, err
	}

	return driveFile, nil
}

func syncFolder(svc *drive.Service, folderName string) error {
	dirPath := constants.WorkDir + "/" + folderName

	var parentId string
	b, err := ioutil.ReadFile(dirPath + "/" + constants.DriveIdFileName)
	if err != nil {
		driveDir, err := createDir(svc, folderName, constants.DriveBaseFolderId)
		if err != nil {
			log.Error.Println(err)
			return err
		}

		idFile, err := os.OpenFile(dirPath+"/"+constants.DriveIdFileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			log.Error.Println(err)
			return err
		}
		_, err = idFile.WriteString(driveDir.Id)
		if err != nil {
			log.Error.Println(err)
			return err
		}

		parentId = driveDir.Id
	} else {
		parentId = string(b)
	}

	_, err = updateLogFile(svc, dirPath, parentId)
	if err != nil {
		log.Error.Println(err)
		return err
	}

	localFiles, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Error.Println(err)
		return err
	}

	for _, localFile := range localFiles {
		fileName := localFile.Name()
		if fileName[len(fileName)-4:] == ".png" {
			file, err := os.Open(dirPath + "/" + fileName)
			if err != nil {
				log.Error.Println(err)
				return err
			}

			_, err = uploadImage(svc, fileName, file, parentId)
			if err != nil {
				log.Error.Println(err)
				return err
			}

			file.Close()

			err = os.Remove(dirPath + "/" + fileName)
			if err != nil {
				log.Error.Println(err)
				return err
			}
		}
	}

	return nil
}

func Sync() {
	svc, err := getService()
	if err != nil {
		log.Error.Println(err)
		return
	}

	folders, err := ioutil.ReadDir(constants.WorkDir)
	if err != nil {
		log.Error.Println(err)
		return
	}

	for _, folder := range folders {
		if folder.IsDir() {
			folderTime, err := time.Parse("20060102", folder.Name())
			if err != nil {
				log.Warn.Printf("not a log directory: %v", err)
				continue
			}

			err = syncFolder(svc, folder.Name())
			if err != nil {
				log.Error.Println(err)
			}

			if time.Now().Sub(folderTime).Hours() > 170 {
				err = os.RemoveAll(constants.WorkDir + "/" + folder.Name())
				if err != nil {
					log.Error.Println(err)
					return
				}
			}
		}
	}
}
