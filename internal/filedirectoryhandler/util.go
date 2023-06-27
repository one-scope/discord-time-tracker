package filedirectoryhandler

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

// 外で定義する必要あるかな？
const (
	usersDirectoryPath    = "users"
	usersNewJsonFileName  = "newusers.json"
	usersJsonFileName     = "users.json"
	statusesDirectoryPath = "statuses"
	statusesNewFileName   = "newstatuses.json"
	statusesFileName      = "statuses.json"
)

func CreateDataDirectory(aDataBasePath string) error {
	if tError := os.MkdirAll(filepath.Join(aDataBasePath, usersDirectoryPath), 0777); tError != nil {
		return tError
	}
	return os.MkdirAll(filepath.Join(aDataBasePath, statusesDirectoryPath), 0777)
}

func RenameUserFile(aDataBasePath string) error {
	return os.Rename(filepath.Join(aDataBasePath, usersDirectoryPath, usersNewJsonFileName), filepath.Join(aDataBasePath, usersDirectoryPath, usersJsonFileName))
}

func IsExistsUsersFile(aDataBasePath string) bool {
	_, tError := os.Stat(filepath.Join(aDataBasePath, usersDirectoryPath, usersJsonFileName))
	return !os.IsNotExist(tError)
}
func DecodeUsersFile(aDataBasePath string, aData any) error {
	tByte, tError := ioutil.ReadFile(filepath.Join(aDataBasePath, usersDirectoryPath, usersJsonFileName))
	if tError != nil {
		return tError
	}
	return json.Unmarshal(tByte, aData)
}
func EncodeUsersFile(aDataBasePath string, aData any) error {
	tByte, tError := json.MarshalIndent(aData, "", "  ")
	if tError != nil {
		return tError
	}
	return ioutil.WriteFile(filepath.Join(aDataBasePath, usersDirectoryPath, usersNewJsonFileName), tByte, 0777)
}
