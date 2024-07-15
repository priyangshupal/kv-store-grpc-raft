package encoding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncoderDecoder(t *testing.T) {
	test_text := "1234 text"

	encodedText, _ := Encode(test_text)
	decodedText, _ := Decode(encodedText)

	assert.Equal(t, test_text, decodedText, "text should remain unchanged after encoding/decoding")
}
