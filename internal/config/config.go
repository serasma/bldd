package config

import "strings"

type Config struct {
	Libraries   []string
	Directories []string
	ReportPath  string
	Workers     int
}

func New(libraries, directories, reportPath string, workers int) *Config {
	var splitLibraries []string
	if libraries != "" {
		splitLibraries = strings.Split(libraries, ",")
	}

	var splitDirectories []string
	if directories != "" {
		splitDirectories = strings.Split(directories, ",")
	}

	return &Config{
		Libraries:   splitLibraries,
		Directories: splitDirectories,
		ReportPath:  reportPath,
		Workers:     workers,
	}
}
