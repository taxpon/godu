package godu

import (
	"github.com/cloudfoundry/bytefmt"
	"os"
	"path/filepath"
)

// PrintableRecord is an interface for print du records
type PrintableRecord interface {
	HSize() string
	Size() int64
	RelativePath() string
}

// File is abstracted file object
type File struct {
	path    string
	dirSize int64
	info    os.FileInfo
}

// NewFile creates a new File object with given information
func NewFile(path string, fileInfo os.FileInfo) *File {
	return &File{path, 0, fileInfo}
}

// Name returns the file name of File
func (f *File) Name() string {
	return f.info.Name()
}

// RelativePath returns the relative path of File
func (f *File) RelativePath() string {
	return filepath.Join(f.path, f.info.Name())
}

// Size returns file size of File
func (f *File) Size() int64 {
	if f.IsDir() {
		return f.dirSize
	}
	return f.info.Size()
}

// HSize returns file size of File with human-readable format
func (f *File) HSize() string {
	return bytefmt.ByteSize(uint64(f.Size()))
}

// DirSize returns directory size if File is directory else 0
func (f *File) DirSize() int64 {
	return f.dirSize
}

// SetDirSize set DirSize to File
func (f *File) SetDirSize(size int64) {
	f.dirSize = size
}

// IsDir returns whether File is directory or not
func (f *File) IsDir() bool {
	return f.info.IsDir()
}
