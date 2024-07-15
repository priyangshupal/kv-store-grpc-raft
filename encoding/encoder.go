package encoding

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

func Encode(data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(data)
	if err != nil {
		return nil, fmt.Errorf("error encoding data: %v", err)
	}
	return buf.Bytes(), nil
}
