package decompress

import (
	"io"
	"os"

	"github.com/klauspost/compress/zstd"
)

func Decompress(in io.Reader, out io.Writer) error {
	decoder, err := zstd.NewReader(in)
	if err != nil {
		return err
	}
	defer decoder.Close()

	_, err = io.Copy(out, decoder)
	return err
}

func CanDecompress(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer file.Close()

	decoder, err := zstd.NewReader(file)
	if err != nil {
		return false, err
	}
	defer decoder.Close()

	bytes := make([]byte, 100)
	_, err = decoder.Read(bytes)

	if err != nil {
		return false, err
	}

	return true, nil
}
