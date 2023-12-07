package main

import (
	"encoding/json"
	"os"
	"strings"
)

func LoadSettings(filename string) (string, error) {
	filename, err := ExpandHomeDir(filename)
	if err != nil {
		return filename, err
	}
	file, err := os.ReadFile(filename)
	if err != nil {
		return filename, err
	}
	err = json.Unmarshal(file, &filename)
	return filename, err
}

func SaveSettings(settings string, filename string) error {
	file, err := json.MarshalIndent(settings, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, file, 0644)
}

func ExpandHomeDir(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return strings.Replace(path, "~", home, 1), nil
	}
	return path, nil
}
