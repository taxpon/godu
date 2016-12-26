package godu

import (
	"github.com/cloudfoundry/bytefmt"
	"time"
)

// DumpHeader describes general information of dumped records
type DumpHeader struct {
	Version   string
	Path      string
	TimeStamp string
}

// DumpRecord describes one file information
type DumpRecord struct {
	Path    string
	RawSize int64
}

// DumpFile holes single header and multiple records
type DumpFile struct {
	Header  DumpHeader
	Records []DumpRecord
}

// NewDumpHeader creates new DumpHeader instance (pointer)
func NewDumpHeader(path string) *DumpHeader {
	d := &DumpHeader{
		Version:   VERSION,
		Path:      path,
		TimeStamp: time.Now().UTC().Format("2006-01-02T15:04:05-0700"),
	}
	return d
}

// MakeDumpHeader creates new DumpHeader instance (value)
func MakeDumpHeader(path string) DumpHeader {
	return *NewDumpHeader(path)
}

// NewDumpRecord creates new DumpRecord instance (pointer)
func NewDumpRecord(f *File) *DumpRecord {
	d := &DumpRecord{
		Path:    f.RelativePath(),
		RawSize: f.Size(),
	}
	return d
}

// MakeDumpRecord creates new DumpRecord instance (value)
func MakeDumpRecord(f *File) DumpRecord {
	return *NewDumpRecord(f)
}

// HSize returns file size of File with human-readable format
func (d DumpRecord) HSize() string {
	return bytefmt.ByteSize(uint64(d.Size()))
}

// Size returns file size of File
func (d DumpRecord) Size() int64 {
	return d.RawSize
}

// RelativePath returns the relative path of File
func (d DumpRecord) RelativePath() string {
	return d.Path
}
