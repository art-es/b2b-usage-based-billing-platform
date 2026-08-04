package log

import (
	"bufio"
	"bytes"
)

type Buffer struct {
	*bytes.Buffer
}

func NewBuffer() Buffer {
	return Buffer{
		Buffer: &bytes.Buffer{},
	}
}

func (b *Buffer) Logs() []string {
	var lines []string

	s := bufio.NewScanner(b)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		lines = append(lines, s.Text())
	}

	return lines
}
