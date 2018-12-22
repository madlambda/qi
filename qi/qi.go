package qi

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type (
	Filesystem interface {
		Open(name string) (File, error)
		Stat(name string) (os.FileInfo, error)
	}

	// File in the filesystem
	File interface {
		Name() string
		Perm() uint32

		io.Closer
		io.Reader
		//io.ReaderAt
		//io.Seeker
		//Stat() (os.FileInfo, error)
	}

	// MemFS is an in-memory filesystem
	MemFS struct {
		Root *MemDir
	}

	// MemFile is an in-memory file
	MemFile struct {
		name string
		perm uint32

		data []byte
		pos  int
	}

	// MemDir is an in-memory directory
	MemDir struct {
		name string
		perm uint32

		entries []File
	}

	Walkfn func(f File)
)

// NewRoot builds the chaos root filesystem
func NewRoot() *MemFS {
	fs := NewMemFS()

	root := NewMemDir("/", 0644)
	dev := NewMemDir("dev", 0644)

	root.AddFile(dev)
	fs.Root = root
	return fs
}

func (fs *MemFS) Append(dirname string, f File) error {
	dir, err := fs.GetDir(dirname)
	if err != nil {
		return err
	}

	dir.AddFile(f)
	return nil
}

// GetDir gets a directory named path. Requires absolute path
func (fs *MemFS) GetDir(path string) (*MemDir, error) {
	if path == "/" {
		return fs.Root, nil
	}

	pathParts := strings.Split(path, "/")
	idx := 1
	dir := fs.Root

outer:
	for idx < len(pathParts) {
		for _, entry := range dir.entries {
			if entry.Name() != pathParts[idx] {
				continue
			}

			if entry.Perm()&uint32(os.ModeDir) != uint32(os.ModeDir) {
				return nil, fmt.Errorf("file not found 1")
			}

			dir = entry.(*MemDir)
			idx++
			goto outer
		}

		return nil, fmt.Errorf("file not found (%d, %d)", idx, len(pathParts))
	}

	return dir, nil
}

func (fs *MemFS) Walk(path string, f Walkfn) error {
	dir, err := fs.GetDir(path)
	if err != nil {
		return err
	}

	for _, entry := range dir.entries {
		f(entry)
	}

	return nil
}

func NewMemFS() *MemFS {
	return &MemFS{}
}

// NewFile creates a new in-memory file.
func NewFile(name string, perm uint32) *MemFile {
	return &MemFile{
		name: name,
		perm: perm,
	}
}

func (f *MemFile) Name() string { return f.name }
func (f *MemFile) Perm() uint32 { return f.perm }

// NewMemDir creates a new in-memory directory.
func NewMemDir(name string, perm uint32) *MemDir {
	return &MemDir{
		name:    name,
		perm:    perm | uint32(os.ModeDir),
		entries: nil,
	}
}

func (d *MemDir) Name() string { return d.name }
func (d *MemDir) Perm() uint32 { return d.perm }

func (d *MemDir) Read(data []byte) (int, error) {
	// TODO(i4k): should we return a string representation of file entries?
	return 0, fmt.Errorf("read at directory")
}

func (d *MemDir) Write(data []byte) (int, error) {
	return 0, fmt.Errorf("write at directory")
}

// Close the directory
func (d *MemDir) Close() error { return fmt.Errorf("not implemented") }

// AddFile adds file f into directory d.
func (d *MemDir) AddFile(f File) {
	d.entries = append(d.entries, f)
}

func (f *MemFile) Read(data []byte) (int, error) {
	return 0, fmt.Errorf("not implemented")
}

func (f *MemFile) Write(data []byte) (int, error) {
	return 0, fmt.Errorf("not implemented")
}

func (f *MemFile) Close() error {
	return fmt.Errorf("not implemented")
}
