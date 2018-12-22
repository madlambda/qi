package dev

import (
	"fmt"
	"io"
	"sync"
)

type (
	mouseStatus struct {
		x, y int
	}

	mouseFile struct {
		updates chan mouseStatus

		mu  sync.RWMutex
		pos int64
	}
)

const statusLen = 11 + 1 + 11 // 00000000000 00000000000

// MouseInit initializes the mouse driver
func MouseInit() *mouseFile {
	m := &mouseFile{
		updates: make(chan mouseStatus, 1),
	}
	return m
}

func (m *mouseFile) Name() string { return "mouse" }
func (m *mouseFile) Perm() uint32 { return 0644 }

func (m *mouseFile) UpdateCoords(x, y int) {
	m.updates <- mouseStatus{x, y}
}

func (m *mouseFile) Read(data []byte) (int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	status := <-m.updates

	content := fmt.Sprintf("%011d %011d", status.x, status.y)
	if m.pos >= int64(len(content)) {
		return 0, io.EOF
	}

	n := copy(data[0:], []byte(content)[m.pos:])
	m.pos += int64(n)
	return n, nil
}

func (m *mouseFile) Write(data []byte) (int, error) {
	return 0, fmt.Errorf("permission denied")
}

func (m *mouseFile) Seek(offset int64, whence int) (int64, error) {
	if whence == io.SeekStart {
		if offset < 0 {
			return 0, fmt.Errorf("seek: invalid offset")
		}

		m.pos = offset
		return m.pos, nil
	}

	if whence == io.SeekCurrent {
		m.pos += offset
		return m.pos, nil
	}

	if whence == io.SeekEnd {
		m.pos = statusLen - offset
		return m.pos, nil
	}

	return 0, fmt.Errorf("invalid seek whence %d", whence)
}

func (m *mouseFile) Close() error {
	return nil
}
