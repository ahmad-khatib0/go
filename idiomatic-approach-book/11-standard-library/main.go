package main

import (
	"compress/gzip"
	"io"
	"os"
)

func main() {

}

// ********************************* io and Friends *********************************
type NotHowReaderIsDefined interface {
	Read() (p []byte, err error)
}

func countLetters(r io.Reader) (map[string]int, error) {
	buf := make([]byte, 2048)
	out := map[string]int{}
	for {
		n, err := r.Read(buf)
		for _, b := range buf[:n] { // [:n]  => to know how many bytes were written to the buffer.
			// iterate over a sub-slice of our buf slice, processing the data that was read
			if (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') {
				out[string(b)]++
			}
		}
		if err == io.EOF {
			return out, nil
		}
		if err != nil {
			return nil, err
		}
	}
}

func buildGZipReader(fileName string) (*gzip.Reader, func(), error) {
	// This function demonstrates the way to properly wrap types that implement io.Reader.
	r, err := os.Open(fileName)
	if err != nil {
		return nil, nil, err
	}
	gr, err := gzip.NewReader(r)
	if err != nil {
		return nil, nil, err
	}
	return gr, func() {
		gr.Close()
		r.Close()
	}, nil
}
