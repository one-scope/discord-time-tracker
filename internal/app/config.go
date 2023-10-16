package app

import (
	"os"
	"path"
)

type LogConfig struct {
	FilePath string
}

func LoadLogConfig() LogConfig {
	return LogConfig{
		FilePath: path.Join(os.Getenv("LOG_FILE_BASE_PATH"), os.Getenv("LOG_FILE_NAME")),
	}
}
