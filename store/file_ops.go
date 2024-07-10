package store

import (
	"fmt"
	"io"
	"os"
)

// opens the file for file operations or creates
// the file if it doesn't exist
func openFile(fileName string) (*os.File, error) {
	dirPath := "../res"
	// create directory for the file
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return nil, err
	}
	filePath := fmt.Sprintf("%s/%s.txt", dirPath, fileName)

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		return nil, err
	}
	return file, nil
}
func appendToFile(fileName string, data string) error {
	file, err := openFile(fileName)
	if err != nil {
		return err
	}
	_, err = file.WriteString(data)
	if err != nil {
		return err
	}
	file.Close()
	return nil
}
func readFile(fileName string) (string, error) {
	file, err := openFile(fileName)
	if err != nil {
		return "", err
	}
	// Move the file pointer to the beginning of the file
	if _, err = file.Seek(0, 0); err != nil {
		return "", err
	}
	// Read the file contents
	data, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	file.Close()
	return string(data), nil
}
func clearFileContents(fileName string) error {
	filePath := fmt.Sprintf("../res/%s.txt", fileName)
	_, err := os.Stat(filePath)
	if err != nil {
		return err // file doesn't exist
	}
	// Open the file with the O_TRUNC flag to clear its contents
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	return file.Close()
}
func deleteDir(dirPath string) error {
	return os.RemoveAll(dirPath)
}
