package fileops

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadAndWriteFile(t *testing.T) {
	const (
		fileName = "test.txt"
		dirPath  = "res"
	)
	fileContent := []byte("Test content")

	err := WriteToFile(dirPath, fileName, fileContent)
	assert.Equal(t, nil, err, "no error should be present while writing to file")

	content, err := ReadFile(dirPath, fileName)
	assert.Equal(t, nil, err, "no error should be present while reading content from file")
	assert.Equal(t, fileContent, content, "file content should be same as content written")

	deleteDir(dirPath)
}
