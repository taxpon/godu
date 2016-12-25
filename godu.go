package godu

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	msgpack "gopkg.in/vmihailenco/msgpack.v2"
	"io/ioutil"
)


// Run runs godu main logic
func Run(path string, recursive bool, absolute bool, dumpFlg bool, loadFlg bool) error {
	items, size := Scan(path, recursive)

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	cwdAbsPath, err := filepath.Abs(cwd)
	fmt.Println(cwdAbsPath)
	if err != nil {
		return err
	}

	if dumpFlg {
		dumpRecords(cwdAbsPath, "/Users/takuro/Desktop/dump.bin", items)
		return nil
	}

	if loadFlg {
		_, err := loadRecords("/Users/takuro/Desktop/dump.bin")
		if err != nil {
			return err
		}
	}

	for _, item := range items {
		w := color.New(color.FgWhite)
		w.Printf("%6v   ", item.HSize())
		if absolute {
			c := color.New(color.FgCyan)
			c.Printf("%v\n", filepath.Join(cwdAbsPath, item.RelativePath()))
		} else {
			c := color.New(color.FgCyan)
			c.Printf("%v\n", item.RelativePath())
		}

	}
	fmt.Printf("Total size is %v byte.\n", size)
	return nil
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

		if file.IsDir() {

			var subFiles [](*File)
			var subSize int64

			if recursive {
				subFiles, subSize = Scan(file.RelativePath(), recursive)
				*fileList = append(*fileList, subFiles...)
			} else {
				_, subSize = Scan(file.RelativePath(), !recursive)
			}

			file.SetDirSize(subSize)
			totalSize += subSize
		} else {
			totalSize += file.Size()
		}
	}
	return *fileList, totalSize
}



func dumpRecords(path, dumpFileName string, files [](*File)) error {

	dFile := &DumpFile{
		Header: *&DumpHeader{
			Version: VERSION,
			Path: path,
		},
		Records: make([]DumpRecord, len(files)),
	}

	for i, file := range files {
		dRec, err := NewDumpRecord(file)
		if err != nil {
			return err
		}
		dFile.Records[i] = *dRec
	}

	//fmt.Println(dFile.Header)
	//fmt.Println(dFile.Records)

	b, err := msgpack.Marshal(dFile)
	if err != nil {
		return err
	}

	err2 := ioutil.WriteFile(dumpFileName, b, 0644)
	if err2 != nil {
		return err2
	}

	var dFile2 DumpFile
	//dRecs2 := make([]DumpRecord, len(files))
	msgpack.Unmarshal(b, &dFile2)
	fmt.Println(dFile2)

	return nil
}

func loadRecords(dumpFileName string) (*DumpFile, error) {
	rawRecords, err := ioutil.ReadFile(dumpFileName)
	if err != nil {
		return nil, err
	}

	var dFile DumpFile

	err2 := msgpack.Unmarshal(rawRecords, &dFile)
	if err2 != nil {
		return nil, err2
	}

	fmt.Println(dFile)
	return &dFile, nil
}
