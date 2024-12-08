package run_8

import (
	"bytes"
	"io"
	"os"
	"unsafe"
)

// ==================================================================================== //
// StationScanner
// ==================================================================================== //

type StationScanner struct {
	f *os.File

	chunk []byte
	start int
	end   int

	eof bool
}

func newStationScanner(f *os.File) *StationScanner {
	return &StationScanner{
		f:     f,
		chunk: make([]byte, chunkSize),
	}
}

// intTemp converts the second part of a line to int.
// "-77.7" => -777
// "77.7" => 777
func (s *StationScanner) intTemp(bs []byte) int {
	neg := bs[0] == '-'
	if neg {
		bs = bs[1:]
	}

	var res int
	if bs[1] == '.' {
		res = int(bs[0]-48)*10 + int(bs[2]-48)
	} else { // bs[2] == '.'
		res = int(bs[0]-48)*100 + int(bs[1]-48)*10 + int(bs[3]-48)
	}

	if neg {
		return -res
	}
	return res
}

func (s *StationScanner) updateChunk() {
	if s.end-s.start >= maxLineLength && !s.eof {
		return // still at least one whole line left in s.chunk
	}

	// backup necessary to be able to use unsafe in [Line]
	backup := s.chunk
	s.chunk = make([]byte, chunkSize)
	copy(s.chunk, backup)

	copy(s.chunk[:], s.chunk[s.start:s.end])
	s.end = s.end - s.start
	s.start = 0

	n, err := s.f.Read(s.chunk[s.end:])
	if err == io.EOF {
		s.eof = true
	}
	if err != nil && err != io.EOF {
		panic(err)
	}
	s.end += n
}

func (s *StationScanner) Next() bool {
	s.updateChunk()
	return !s.eof || s.start < s.end
}

// Line takes the current chunk and processes the first line in it.
// s.start is advanced by the bytes used/processed.
func (s *StationScanner) Line() (name string, temp int) {
	lines := s.chunk[s.start:]

	l := bytes.IndexByte(lines, ';')

	bName := lines[:l]
	name = unsafe.String(unsafe.SliceData(bName), len(bName))

	switch {
	case lines[l+4] == '\n': // 1.2
		s.start += l + 5 // increment the start position by the bytes used
		return name, s.intTemp(lines[l+1 : l+4])
	case lines[l+5] == '\n': // 12.3 or -1.2
		s.start += l + 6
		return name, s.intTemp(lines[l+1 : l+5])
	case lines[l+6] == '\n': // -12.3
		s.start += l + 7
		return name, s.intTemp(lines[l+1 : l+6])
	default:
		panic("not a line")
	}
}
