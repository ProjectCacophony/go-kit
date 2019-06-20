package events

import (
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jinzhu/gorm"
)

const (
	noStorageError string = "Event storage is Nil"
	noFileData     string = "No file data."
	noFileInfo     string = "No file info found."

	fileIdCharacters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	fileIdLength     = 8
)

var (
	objectStorageFQDN       = ""
	objectStorageBucketName = ""
)

type FileInfo struct {
	gorm.Model

	Filename        string
	UserID          string
	ChannelID       string
	GuildID         string
	FileID          string
	Source          string
	MimeType        string
	UploadDate      time.Time
	Filesize        int
	RetrievedCount  int
	Public          bool
	CustomCommandID uint
}

func (f *FileInfo) TableName() string {
	return "object_storage"
}

func (f *FileInfo) GetLink() string {
	if objectStorageFQDN != "" {
		return "https://" + objectStorageFQDN + "/" + f.bucketKey()
	}

	if objectStorageBucketName != "" {
		return "https://storage.googleapis.com/" + objectStorageBucketName + "/" + f.bucketKey()
	}

	return ""
}

func (f *FileInfo) bucketKey() string {
	return f.FileID + "/" + f.Filename
}

func InitObjectStorage(db *gorm.DB, fqdn string, bucketName string) error {

	// objectStorageFQDN = fqdn
	objectStorageBucketName = bucketName

	rand.Seed(time.Now().UTC().UnixNano())
	return db.AutoMigrate(FileInfo{}).Error
}

func (e *Event) AddFile(data []byte, file *FileInfo) (*FileInfo, error) {
	if e.storage == nil {
		return nil, errors.New(noStorageError)
	}

	if len(data) == 0 {
		return nil, errors.New(noFileData)
	}

	if file.FileID == "" {
		newFileID, err := getUniqueFileId(e.DB())
		if err != nil {
			return nil, err
		}
		file.FileID = newFileID
	}

	err := saveFileToDB(e.DB(), file)
	if err != nil {
		return nil, err
	}

	err = e.storage.bucket.WriteAll(e.Context(), file.bucketKey(), data, nil)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (e *Event) DeleteFile(file *FileInfo) error {
	if e.storage == nil {
		return errors.New(noStorageError)
	}

	if file.Filename == "" {
		return errors.New(noFileInfo)
	}

	err := deleteFileFromDB(e.DB(), file)
	if err != nil {
		return err
	}

	return e.storage.bucket.Delete(e.Context(), file.bucketKey())
}

func (e *Event) AddAttachement(attachement *discordgo.MessageAttachment) (*FileInfo, error) {

	resp, err := e.HTTPClient().Get(attachement.URL)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	uniqueId, err := getUniqueFileId(e.DB())
	if err != nil {
		return nil, err
	}

	return e.AddFile(bytes, &FileInfo{
		Filename:   attachement.Filename,
		UserID:     e.UserID,
		ChannelID:  e.ChannelID,
		GuildID:    e.GuildID,
		FileID:     uniqueId,
		MimeType:   http.DetectContentType(bytes),
		UploadDate: time.Now(),
		Source:     attachement.URL,
		Filesize:   len(bytes),
		Public:     true,
	})

}

func (e *Event) UpdateFileInfo(file FileInfo) error {
	return e.DB().Update(file).Error
}

func getUniqueFileId(db *gorm.DB) (string, error) {

	output := make([]byte, fileIdLength)
	for i := range output {
		output[i] = fileIdCharacters[rand.Intn(len(fileIdCharacters))]
	}

	var count int
	err := db.
		Model(FileInfo{}).
		Where(FileInfo{FileID: string(output)}).
		Count(&count).
		Error
	if count != 0 {
		return getUniqueFileId(db)
	}

	return string(output), err
}

func saveFileToDB(db *gorm.DB, f *FileInfo) error {
	if f.FileID == "" {
		return errors.New(noFileInfo)
	}

	if db.NewRecord(f) {
		return db.Create(f).Error
	}
	return db.Update(f).Error
}

func deleteFileFromDB(db *gorm.DB, f *FileInfo) error {
	if f.FileID == "" {
		return errors.New(noFileInfo)
	}
	return db.Delete(f).Error
}
