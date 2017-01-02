package godu_test

import (
	"testing"
	"github.com/taxpon/godu"
)

func TestCompareMap_GetValue(t *testing.T) {

}


func TestCompareMap_CompareNew(t *testing.T) {
	c1 := &godu.CompareMap{
		Map: map[string]int64{
			"aaa": 124,
		},
	}

	c2 := &godu.CompareMap{
		Map: map[string]int64{
		},
	}

	cr := c1.Compare(c2)
	if len(cr.New) != 1 {
		t.Errorf("Invalid number of new files, expected %d got %v", 1, len(cr.Updated))
	}

	if cr.New[0].Path != "aaa" {
		t.Errorf("Invalid Path of new file, expected %s got %v", "aaa", cr.Updated[0].Path)
	}

	if cr.New[0].BeforeSize != 0 {
		t.Errorf("Invalid BeforeSize of new file, expected %s got %v", 0, cr.Updated[0].BeforeSize)
	}

	if cr.New[0].AfterSize != 124 {
		t.Errorf("Invalid AfterSize of new file, expected %s got %v", 124, cr.Updated[0].AfterSize)
	}
}

func TestCompareMap_CompareUpdated(t *testing.T) {
	c1 := &godu.CompareMap{
		Map: map[string]int64{
			"aaa": 124,
		},
	}

	c2 := &godu.CompareMap{
		Map: map[string]int64{
			"aaa": 125,
		},
	}

	cr := c1.Compare(c2)
	if len(cr.Updated) != 1 {
		t.Errorf("Invalid number of updated files, expected %d got %v", 1, len(cr.Updated))
	}

	if cr.Updated[0].Path != "aaa" {
		t.Errorf("Invalid Path of updated file, expected %s got %v", "aaa", cr.Updated[0].Path)
	}

	if cr.Updated[0].BeforeSize != 125 {
		t.Errorf("Invalid BeforeSize of updated file, expected %s got %v", 125, cr.Updated[0].BeforeSize)
	}

	if cr.Updated[0].AfterSize != 124 {
		t.Errorf("Invalid AfterSize of updated file, expected %s got %v", 124, cr.Updated[0].AfterSize)
	}
}

func TestCompareMap_CompareDeleted(t *testing.T) {
	c1 := &godu.CompareMap{
		Map: map[string]int64{
		},
	}

	c2 := &godu.CompareMap{
		Map: map[string]int64{
			"aaa": 125,
		},
	}

	cr := c1.Compare(c2)
	if len(cr.Deleted) != 1 {
		t.Errorf("Invalid number of deleted files, expected %d got %v", 1, len(cr.Deleted))
	}

	if cr.Deleted[0].Path != "aaa" {
		t.Errorf("Invalid Path of deleted file, expected %s got %v", "aaa", cr.Deleted[0].Path)
	}

	if cr.Deleted[0].BeforeSize != 125 {
		t.Errorf("Invalid BeforeSize of deleted file, expected %s got %v", 125, cr.Deleted[0].BeforeSize)
	}

	if cr.Deleted[0].AfterSize != 0 {
		t.Errorf("Invalid AfterSize of deleted file, expected %s got %v", 124, cr.Deleted[0].AfterSize)
	}
}