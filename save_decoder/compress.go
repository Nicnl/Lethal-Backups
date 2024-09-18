package save_decoder

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
)

func base64LineBreaks(input string, length int) string {
	var buf bytes.Buffer
	for i := 0; i < len(input); i += length {
		end := i + length
		if end > len(input) {
			end = len(input)
		}
		buf.WriteString(input[i:end])
		buf.WriteByte('\n')
	}
	return buf.String()
}

func CompressSave(jsonSave JsonSave) (string, error) {
	buf := bytes.Buffer{}
	gzipWriter, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return "", err
	}

	_, err = gzipWriter.Write(jsonSave.Data)
	if err != nil {
		return "", err
	}

	err = gzipWriter.Close()
	if err != nil {
		return "", err
	}

	// to base64 string with line breaks
	base64String := base64.StdEncoding.EncodeToString(buf.Bytes())
	base64String = base64LineBreaks(base64String, 128)

	return base64String, nil
}

func DecompressSave(base64Save string) (JsonSave, error) {
	// remove line breaks
	gzipSave, err := base64.StdEncoding.DecodeString(base64Save)
	if err != nil {
		return JsonSave{}, err
	}

	gzipReader, err := gzip.NewReader(bytes.NewBuffer(gzipSave))
	if err != nil {
		return JsonSave{}, err
	}
	defer gzipReader.Close()

	var decompressed bytes.Buffer
	_, err = decompressed.ReadFrom(gzipReader)
	if err != nil {
		return JsonSave{}, err
	}

	return JsonSave{decompressed.Bytes()}, nil
}
