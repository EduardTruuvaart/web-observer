package compressor

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompressThenDataIsCompressedCorrectly(t *testing.T) {
	// Arrange
	data := []byte("<html><body><span>Out of Stock</span>My content</body></html>")

	// Act
	compressedData, err := Compress(data)
	encodedData := base64.StdEncoding.EncodeToString(compressedData)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "H4sIAAAAAAAA/7LJKMnNsbNJyk+ptLMpLkjMs/MvLVHIT1MILslPzrbRBwv5Viok5+eVpOaV2OhDVOqDtQECAAD//1Jq3a09AAAA", encodedData)
}

func TestDecompressThenRestoresCorrectly(t *testing.T) {
	// Arrange
	data, _ := base64.StdEncoding.Strict().DecodeString("H4sIAAAAAAAA/7LJKMnNsbNJyk+ptLMpLkjMs/MvLVHIT1MILslPzrbRBwv5Viok5+eVpOaV2OhDVOqDtQECAAD//1Jq3a09AAAA")

	// Act
	decompressedData, err := Decompress(data)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "<html><body><span>Out of Stock</span>My content</body></html>", string(decompressedData))

}
