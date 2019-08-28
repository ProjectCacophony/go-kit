package events

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"mime"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jinzhu/gorm"
	"gitlab.com/Cacophony/go-kit/permissions"
)

const (
	NoStorageSpace          string = "common.noStorageSpace"
	NoStoragePermission     string = "common.noStoragePermission"
	FileTooBig              string = "common.fileTooBig"
	CouldNotExtractFilename string = "common.couldNotExtractFilename"

	noStorageError string = "Event storage is Nil"
	noFileData     string = "No file data."
	noFileInfo     string = "No file info found."

	fileIDCharacters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	fileIDLength     = 8

	userStorageLimit = 1000000000 //   1 GB
	maxUploadLimit   = 100000000  // 100 MB
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

type UserStorageInfo struct {
	FileCount        int
	StorageUsed      int
	StorageAvailable int
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

	objectStorageFQDN = fqdn
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

	usageInfo, err := e.GetUserStorageUsage()
	if err != nil {
		return nil, err
	}

	if !e.Has(permissions.CacoFileStorage) {
		return nil, errors.New(NoStoragePermission)
	}

	if usageInfo.StorageAvailable >= 0 &&
		(usageInfo.StorageUsed+file.Filesize) > usageInfo.StorageAvailable {
		return nil, errors.New(NoStorageSpace)
	}

	if file.FileID == "" {
		newFileID, err := getUniqueFileID(e.DB())
		if err != nil {
			return nil, err
		}
		file.FileID = newFileID
	}

	err = saveFileToDB(e.DB(), file)
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

func (e *Event) AddFileFromURL(link string, filename string) (*FileInfo, error) {
	resp, err := e.HTTPClient().Get(link)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if resp.ContentLength > maxUploadLimit {
		return nil, errors.New(FileTooBig)
	}

	limitedReader := &io.LimitedReader{R: resp.Body, N: maxUploadLimit + 1}

	bytes, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		return nil, err
	}

	// check if limited reader is exhausted, and if so stop, because we do not want to store a partial file
	if limitedReader.N <= 0 {
		return nil, errors.New(FileTooBig)
	}

	// try to extract filename if no filename is passed
	if filename == "" {
		filename = path.Base(resp.Request.URL.Path)

		_, params, err := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
		if err == nil && params["filename"] != "" {
			filename = params["filename"]
		}
	}
	if filename == "" || !strings.Contains(filename, ".") {
		return nil, errors.New(CouldNotExtractFilename)
	}

	return e.AddFile(bytes, &FileInfo{
		Filename:   filename,
		UserID:     e.UserID,
		ChannelID:  e.ChannelID,
		GuildID:    e.GuildID,
		MimeType:   http.DetectContentType(bytes),
		UploadDate: time.Now(),
		Source:     link,
		Filesize:   len(bytes),
		Public:     true,
	})
}

func (e *Event) AddAttachement(attachement *discordgo.MessageAttachment) (*FileInfo, error) {
	return e.AddFileFromURL(attachement.URL, attachement.Filename)
}

func (e *Event) UpdateFileInfo(file FileInfo) error {
	return e.DB().Update(file).Error
}

func (e *Event) GetUserStorageUsage() (*UserStorageInfo, error) {

	info := &UserStorageInfo{
		FileCount:        0,
		StorageUsed:      0,
		StorageAvailable: userStorageLimit,
	}

	err := e.DB().
		Table((&FileInfo{}).TableName()).
		Select("count(*) as file_count, sum(filesize) as storage_used").
		Where("user_id = ?", e.UserID).
		Find(&info).Error

	if e.Has(permissions.BotAdmin) {
		info.StorageAvailable = -1
	}

	return info, err
}

func getUniqueFileID(db *gorm.DB) (string, error) {

	output := make([]byte, fileIDLength)
	for i := range output {
		output[i] = fileIDCharacters[rand.Intn(len(fileIDCharacters))]
	}

	var count int
	err := db.
		Model(FileInfo{}).
		Where(FileInfo{FileID: string(output)}).
		Count(&count).
		Error
	if count != 0 {
		return getUniqueFileID(db)
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
