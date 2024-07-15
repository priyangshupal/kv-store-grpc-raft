package fileops

import (
	"fmt"
	"io"
	"os"
)

const FILE_EXTENSION = "bin"

// opens the file for file operations or creates
// the file if it doesn't exist
func openFile(dirPath string, fileName string) (*os.File, error) {
	// create directory for the file
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return nil, err
	}
	filePath := fmt.Sprintf("%s/%s.%s", dirPath, fileName, FILE_EXTENSION)

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		return nil, err
	}
	return file, nil
}
func AppendToFile(dirPath string, fileName string, data []byte) error {
	// create directory for the file
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return err
	}
	filePath := fmt.Sprintf("%s/%s.%s", dirPath, fileName, FILE_EXTENSION)

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	file.Close()
	return nil
}
func ReadFile(dirPath string, fileName string) ([]byte, error) {
	file, err := openFile(dirPath, fileName)
	if err != nil {
		return []byte(""), err
	}
	// Move the file pointer to the beginning of the file
	if _, err = file.Seek(0, 0); err != nil {
		return []byte(""), err
	}
	// Read the file contents
	data, err := io.ReadAll(file)
	if err != nil {
		return []byte(""), err
	}
	file.Close()
	return data, nil
}
func clearFileContents(dirPath string, fileName string) error {
	filePath := fmt.Sprintf("%s/%s.%s", dirPath, fileName, FILE_EXTENSION)
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
func WriteToFile(dirPath string, filename string, content []byte) error {
	// create directory for the file
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return err
	}
	filePath := fmt.Sprintf("%s/%s.%s", dirPath, filename, FILE_EXTENSION)

	// Open the file, create it if it doesn't exist, truncate it if it does
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	// Write the content to the file
	_, err = file.Write(content)
	if err != nil {
		return fmt.Errorf("could not write to file: %v", err)
	}

	return nil
}
func deleteDir(dirPath string) error {
	return os.RemoveAll(dirPath)
}
