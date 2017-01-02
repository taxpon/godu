package godu

import (
	"fmt"
	"path/filepath"
)

type recordType int

const (
	created recordType = iota
	updated
	deleted
)

// CompareResultRecord ...
type CompareResultRecord struct {
	Path       string
	BeforeSize int64
	AfterSize  int64
	RecordType recordType
}

// CompareResult ...
type CompareResult struct {
	New     []CompareResultRecord
	Updated []CompareResultRecord
	Deleted []CompareResultRecord
}

// CompareMap ...
type CompareMap struct {
	Map map[string]int64
}

// BuildCompareMap ...
func (f *DumpFile) BuildCompareMap() *CompareMap {
	cm := &CompareMap{Map: make(map[string]int64)}
	for _, record := range f.Records {
		cm.Map[filepath.Join(f.Header.Path, record.Path)] = record.RawSize
	}
	return cm
}

// Compare ...
func (c *CompareMap) Compare(other *CompareMap) *CompareResult {
	cr := &CompareResult{}
	for k, v := range c.Map {
		ov, ok := other.Map[k]
		if !ok {
			cr.addNewFile(k, v)
			continue
		}
		if v != ov {
			cr.addUpdatedFile(k, ov, v)
		}
		delete(other.Map, k)
	}
	for k, v := range other.Map {
		cr.addDeletedFile(k, v)
	}
	return cr
}

func (cr *CompareResult) addNewFile(path string, size int64) {
	cr.New = append(cr.New, CompareResultRecord{
		Path:       path,
		BeforeSize: 0,
		AfterSize:  size,
		RecordType: created,
	})
}

func (cr *CompareResult) addUpdatedFile(path string, beforeSize int64, afterSize int64) {
	cr.Updated = append(cr.Updated, CompareResultRecord{
		Path:       path,
		BeforeSize: beforeSize,
		AfterSize:  afterSize,
		RecordType: updated,
	})
}

func (cr *CompareResult) addDeletedFile(path string, size int64) {
	cr.Deleted = append(cr.Updated, CompareResultRecord{
		Path:       path,
		BeforeSize: size,
		AfterSize:  0,
		RecordType: deleted,
	})
}

// PrintRecord ...
func (crr *CompareResultRecord) PrintRecord() {
	switch crr.RecordType {
	case created:
		fmt.Printf("  %s (+%v)\n", crr.Path, crr.AfterSize)
	case updated:
		fmt.Printf("  %s (%v <- %v)\n", crr.Path, crr.BeforeSize, crr.AfterSize)
	case deleted:
		fmt.Printf("  %s (-%v)\n", crr.Path, crr.BeforeSize)
	}
}
