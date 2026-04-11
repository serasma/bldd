package config

import "strings"

type Config struct {
	Libraries   []string
	Directories []string
	ReportPath  string
	Workers     int
}

func New(libraries, directories, reportPath string, workers int) *Config {
	return &Config{
		Libraries:   strings.Split(libraries, ","),
		Directories: strings.Split(directories, ","),
		ReportPath:  reportPath,
		Workers:     workers,
	}
}
