package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
)

func (u *Utils) ToJSON(src any) []byte {
	if src == nil {
		return nil
	}

	enc, _ := json.Marshal(src)
	return enc
}

func (u *Utils) JsonLineAndCharErr(offset int64, pay []byte) (int, int, error) {
	if offset < 0 {
		return -1, -1, errors.New("offset value cannot be negative")
	}

	br := bytes.NewReader(pay)
	// Count lines and characters.
	lnum := 1
	cnum := 0

	// number of consumed bytes
	var count int64
	for {
		ch, size, err := br.ReadRune()
		if err == io.EOF {
			return -1, -1, errors.New("offset value too large")
		}

		count += int64(size)
		if ch == '\n' {
			lnum++
			cnum = 0
		} else {
			cnum++
		}

		if count >= offset {
			break
		}
	}

	return lnum, cnum, nil
}
