package encoding

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

// Decode bytes back to map[string]string using gob
func DecodeMap(data []byte) (map[string]string, error) {
	var result map[string]string
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("error decoding data: %v", err)
	}
	return result, nil
}

func Decode(encodedData []byte) (string, error) {
	var result string
	buf := bytes.NewBuffer(encodedData)
	decoder := gob.NewDecoder(buf)
	err := decoder.Decode(&result)
	if err != nil {
		return "", err
	}
	return result, nil
}
