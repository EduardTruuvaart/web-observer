package compressor

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"log"
)

func Compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write([]byte(data)); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		log.Fatal(err)
	}

	return b.Bytes(), nil
}

func Decompress(compressedData []byte) ([]byte, error) {
	reader := bytes.NewReader(compressedData)
	gzreader, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}

	output, err := ioutil.ReadAll(gzreader)
	if err != nil {
		return nil, err
	}

	return output, nil
}
