package filedirectoryhandler

import (
	"encoding/json"
	"io/ioutil"
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
func DecodeUsersJsonFile(aDataBasePath string, aData any) error {
	tByte, tError := ioutil.ReadFile(filepath.Join(aDataBasePath, usersDirectoryPath, usersJsonFileName))
	if tError != nil {
		return tError
	}
	return json.Unmarshal(tByte, aData)
}
func EncodeUsersJsonFile(aDataBasePath string, aData any) error {
	tByte, tError := json.MarshalIndent(aData, "", "  ")
	if tError != nil {
		return tError
	}
	return ioutil.WriteFile(filepath.Join(aDataBasePath, usersDirectoryPath, usersNewJsonFileName), tByte, 0777)
}
