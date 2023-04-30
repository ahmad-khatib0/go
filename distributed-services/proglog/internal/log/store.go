package log

import (
	"bufio"
	"encoding/binary"
	"os"
	"sync"
)

var enc = binary.BigEndian // the encoding that we persist record sizes and index entries in
const lenWidth = 8         // the number of bytes used to store the record’s length

type store struct {
	*os.File
	mu   sync.Mutex
	buf  *bufio.Writer
	size uint64
}

func newStore(f *os.File) (*store, error) {

	fi, err := os.Stat(f.Name())
	// ╒═════════════════════════════════════════════════════════════════════════════════════════════════╕
	//   os.Stat(name string) to get the file’s current size, in case we’re re-creating the store from a
	//   file that has existing data, which would happen if, for example, our service had restarted
	// ╘═════════════════════════════════════════════════════════════════════════════════════════════════╛
	if err != nil {
		return nil, err
	}

	size := uint64(fi.Size())

	return &store{
		File: f,
		size: size,
		buf:  bufio.NewWriter(f),
	}, nil
}

func (s *store) Append(p []byte) (n uint64, pos uint64, err error) {

	s.mu.Lock()
	defer s.mu.Unlock()
	pos = s.size

	if err := binary.Write(s.buf, enc, uint64(len(p))); err != nil {
		return 0, 0, err
	}

	//  ╒═══════════════════════════════════════════════════════════════════════════════╕
	//    We write to the buffered writer instead of directly to the file to reduce the
	//    number of system calls and improve performance.
	//  ╘═══════════════════════════════════════════════════════════════════════════════╛
	w, err := s.buf.Write(p) // returns the written bytes number
	if err != nil {
		return 0, 0, err
	}

	//  ╒════════════════════════════════════════════════════════════════════════════════════════════════════╕
	//    We write the length of the record so that, when we read the record, we know how many bytes to read
	//  ╘════════════════════════════════════════════════════════════════════════════════════════════════════╛
	w += lenWidth
	s.size += uint64(w)
	return uint64(w), pos, nil
}

func (s *store) Read(pos uint64) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	//╒═══════════════════════════════════════════════════════════════════╕
	//  it flushes the writer buffer, in case we’re about to try to read
	//  a record that the buffer hasn’t flushed to disk yet.
	//╘═══════════════════════════════════════════════════════════════════╛
	if err := s.buf.Flush(); err != nil {
		return nil, err
	}

	size := make([]byte, lenWidth)
	if _, err := s.File.ReadAt(size, int64(pos)); err != nil {
		return nil, err
	}

	b := make([]byte, enc.Uint64(size))
	if _, err := s.File.ReadAt(b, int64(pos+lenWidth)); err != nil {
		return nil, err
	}

	return b, nil
}

func (s *store) ReadAt(p []byte, off int64) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.buf.Flush(); err != nil {
		return 0, err
	}

	return s.File.ReadAt(p, off)
}

// Close persists any buffered data before closing the file.
func (s *store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.buf.Flush()
	if err != nil {
		return err
	}

	return s.File.Close()
}
