package godu

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/fatih/color"
)


// Run runs godu main logic
func Run(path string, recursive bool, absolute bool) {
	items, size := Scan(path, recursive)

	for _, item := range items {
		w := color.New(color.FgWhite)
		w.Printf("%6v   ", item.HSize())
		if absolute {
			absPath, err := filepath.Abs(path)
			if err != nil {
				fmt.Errorf("Failed to retrieve absokute path for %v", path)
				return
			}
			c := color.New(color.FgCyan)
			c.Printf("%v\n", filepath.Join(absPath, item.FullPath()))
		} else {
			c := color.New(color.FgCyan)
			c.Printf("%v\n", item.FullPath())
		}

	}
	fmt.Printf("Total size is %v byte.\n", size)
}


// Scan scans given directory
func Scan(path string, recursive bool) ([](*File), int64) {
	f, err := os.Open(path)
	if err != nil {
		fmt.Errorf("%s", err)
	}
	defer f.Close()

	fileInfoList, err := f.Readdir(0)
	if err != nil {
		fmt.Errorf("%s", err)
	}

	fileList := &[](*File){}
	totalSize := int64(0)
	for _, fileInfo := range fileInfoList {
		file := NewFile(path, fileInfo)
		*fileList = append(*fileList, file)

		//fmt.Println("addr:", &file)
		//fmt.Printf("File name is %v.\n", file.FullPath())
		//fmt.Printf("FileList pointer is :%p\n", *fileList)

		//for _, _file := range *fileList {
		//	fmt.Printf("  file:%v, pointer:%p\n", _file.Name(), _file)
		//}
		//fmt.Println()

		if file.IsDir() {

			var subFiles [](*File)
			var subSize int64

			if recursive {
				subFiles, subSize = Scan(file.FullPath(), recursive)
				*fileList = append(*fileList, subFiles...)
			} else {
				_, subSize = Scan(file.FullPath(), !recursive)
			}

			file.SetDirSize(subSize)
			totalSize += subSize
		} else {
			totalSize += file.Size()
		}
	}
	return *fileList, totalSize
}
