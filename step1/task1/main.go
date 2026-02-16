package main

import (
	"os"
)

func WriteToLogFile(message string, filePath string) error {
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(message)
	if err != nil {
		return err
	}

	return nil
}
