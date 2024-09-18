package save_historizer

import (
	"bytes"
	"compress/gzip"
)

func compressGzip(input []byte) ([]byte, error) {
	buf := bytes.Buffer{}
	w, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return nil, err
	}

	_, err = w.Write(input)
	if err != nil {
		return nil, err
	}

	err = w.Close()

	return buf.Bytes(), err
}

func decompressGzip(input []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewBuffer(input))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var decompressed bytes.Buffer
	_, err = decompressed.ReadFrom(r)

	return decompressed.Bytes(), err
}
