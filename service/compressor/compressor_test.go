package compressor

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompressThenDataIsCompressedCorrectly(t *testing.T) {
	// Arrange
	data := []byte("<html><body>My content</body></html>")

	// Act
	compressedData, err := Compress(data)
	encodedData := base64.StdEncoding.EncodeToString(compressedData)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, 54, len(compressedData))
	assert.Equal(t, "H4sIAAAAAAAA/7LJKMnNsbNJyk+ptPOtVEjOzytJzSux0QcL2OiDZQEBAAD///fZaEQkAAAA", encodedData)
}

func TestDecompressThenRestoresCorrectly(t *testing.T) {
	// Arrange
	data, _ := base64.StdEncoding.Strict().DecodeString("H4sIAAAAAAAA/7LJKMnNsbNJyk+ptPOtVEjOzytJzSux0QcL2OiDZQEBAAD///fZaEQkAAAA")

	// Act
	decompressedData, err := Decompress(data)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "<html><body>My content</body></html>", string(decompressedData))

}
