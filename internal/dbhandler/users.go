package dbhandler

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
)

const (
	usersDirectoryPath   = "users"
	usersNewJsonFileName = "newusers.json"
	usersJsonFileName    = "users.json"
)

func CreateUsersDataDirectory(aDataBasePath string) error {
	return os.MkdirAll(filepath.Join(aDataBasePath, usersDirectoryPath), 0777)
}

func RenameUsersJsonFile(aDataBasePath string) error {
	return os.Rename(filepath.Join(aDataBasePath, usersDirectoryPath, usersNewJsonFileName), filepath.Join(aDataBasePath, usersDirectoryPath, usersJsonFileName))
}

func IsExistsUsersJsonFile(aDataBasePath string) bool {
	_, tError := os.Stat(filepath.Join(aDataBasePath, usersDirectoryPath, usersJsonFileName))
	return !os.IsNotExist(tError)
}
func DecodeUsersJsonFile(aDataBasePath string, aUserByID *map[string]*User) error {
	tFile, tError := os.OpenFile(filepath.Join(aDataBasePath, usersDirectoryPath, usersJsonFileName), os.O_WRONLY, 0777)
	if tError != nil {
		return tError
	}
	defer func() {
		if tError := tFile.Close(); tError != nil {
			log.Println(tError)
		}
	}()

	tByte, tError := io.ReadAll(tFile)
	if tError != nil {
		return tError
	}
	if tError := json.Unmarshal(tByte, aUserByID); tError != nil {
		return tError
	}

	if aUserByID == nil {
		*aUserByID = map[string]*User{}
	}

	return nil
}

func EncodeUsersJsonFile(aDataBasePath string, aUserByID *map[string]*User) error {
	tFile, tError := os.OpenFile(filepath.Join(aDataBasePath, usersDirectoryPath, usersNewJsonFileName), os.O_WRONLY|os.O_CREATE, 0777)
	if tError != nil {
		return tError
	}
	defer func() {
		if tError := tFile.Close(); tError != nil {
			log.Println(tError)
		}
	}()

	tByte, tError := json.MarshalIndent(aUserByID, "", "  ")
	if tError != nil {
		return tError
	}

	if _, tError := tFile.Write(tByte); tError != nil {
		return tError
	}
	return nil
}
