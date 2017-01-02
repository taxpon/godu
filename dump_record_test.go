package godu_test

import (
	"github.com/taxpon/godu"
	"path/filepath"
	"testing"
)

func TestDumpFile_BuildCompareMap(t *testing.T) {

	cases := [][]map[string]int64{
		{{"a.dat": 100}},
		{{"a.dat": 100}, {"b.dat": 200}},
		{{"a.dat": 100}, {"b.dat": 200}, {"c.dat": 300}},
		{{"a.dat": 100}, {"b.dat": 200}, {"c.dat": 300}, {"d.dat": 400}},
	}

	for _, args := range cases {

		records := []godu.DumpRecord{}
		for _, arg := range args {
			for k, v := range arg {
				records = append(records, godu.DumpRecord{Path: k, RawSize: v})
			}
		}

		df := godu.DumpFile{
			Header:  godu.MakeDumpHeader("/"),
			Records: records,
		}

		cm := df.BuildCompareMap()

		for _, arg := range args {
			for k, v := range arg {
				if cm.Map[filepath.Join("/", k)] != v {
					t.Errorf("expected %v, got %d", v, cm.Map[k])
				}
			}
		}
	}
}
