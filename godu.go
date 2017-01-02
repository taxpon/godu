package godu

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"io/ioutil"

	"github.com/fatih/color"
	msgpack "gopkg.in/vmihailenco/msgpack.v2"
	"time"
)

// Run runs godu main logic
func Run(path string, recursive bool, absolute bool, dumpFlg bool) error {
	items, size := Scan(path, recursive)

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	cwdAbsPath, err := filepath.Abs(cwd)
	if err != nil {
		return err
	}

	if dumpFlg {
		dd, err := getArchivesDir()
		if err != nil {
			return err
		}

		filename := time.Now().UTC().Format("20060102150405.bin")
		err = dumpRecords(cwdAbsPath, filepath.Join(dd, filename), items)
		if err != nil {
			return err
		}
	}

	records := make([]PrintableRecord, len(items))
	for i, record := range items {
		records[i] = record
	}
	if absolute {
		printRecord(records, cwdAbsPath)
	} else {
		printRecord(records, "")
	}

	fmt.Printf("Total size is %v byte.\n", size)
	return nil
}

// Load loads saved file and display the information to stdout
func Load(listFlg bool, filename string) error {
	dd, err := getArchivesDir()
	if err != nil {
		return err
	}

	if listFlg {
		f, err := os.Open(dd)
		if err != nil {
			fmt.Errorf("%s", err)
		}
		defer f.Close()

		fileInfoList, err := f.Readdir(0)
		if err != nil {
			return err
		}

		for _, fileInfo := range fileInfoList {
			fmt.Println(fileInfo.Name())
		}
		return nil
	}

	err = loadAndPrint(filepath.Join(dd, filename))
	return err
}

// Compare compares 2 archived data and show diff result
func Compare(dumpFileName1 string, dumpFileName2 string) error {
	dd, err := getArchivesDir()
	if err != nil {
		return err
	}

	df1, err := loadRecords(filepath.Join(dd, dumpFileName1))
	if err != nil {
		return err
	}

	df2, err := loadRecords(filepath.Join(dd, dumpFileName2))
	if err != nil {
		return err
	}

	cm1 := df1.BuildCompareMap()
	cm2 := df2.BuildCompareMap()
	cr := cm1.Compare(cm2)

	if len(cr.New) > 0 {
		fmt.Println("New:")
		for _, r := range cr.New {
			r.PrintRecord()
		}
		fmt.Println("")
	}

	if len(cr.Updated) > 0 {
		fmt.Println("Updated:")
		for _, r := range cr.Updated {
			r.PrintRecord()
		}
		fmt.Println("")
	}

	if len(cr.Deleted) > 0 {
		fmt.Println("Deleted:")
		for _, r := range cr.Deleted {
			r.PrintRecord()
		}
	}

	return nil
}

func getArchivesDir() (string, error) {
	home := os.Getenv("HOME")
	if home == "" && runtime.GOOS == "windows" {
		home = os.Getenv("APPDATA")
	}

	dir := filepath.Join(home, ".config", "godu", "archives")
	_, err := os.Stat(dir)

	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(dir, 0755)
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}
	return dir, nil
}

func loadAndPrint(dumpFileName string) error {
	dFile, err := loadRecords(dumpFileName)
	if err != nil {
		return err
	}

	records := make([]PrintableRecord, len(dFile.Records))
	for i, record := range dFile.Records {
		records[i] = record
	}
	printHeader(dFile.Header)
	printRecord(records, dFile.Header.Path)
	return nil
}

func printHeader(header DumpHeader) {
	m := color.New(color.FgMagenta)
	m.Printf("godu (version %s)\n", header.Version)
	m.Printf("Saved at %s\n", header.TimeStamp)

	w := color.New(color.FgWhite)
	w.Println("---------------------------------")
}

func printRecord(records []PrintableRecord, absolutePath string) {
	for _, record := range records {
		w := color.New(color.FgWhite)
		w.Printf("%6v   ", record.HSize())
		if absolutePath != "" {
			c := color.New(color.FgCyan)
			c.Printf("%v\n", filepath.Join(absolutePath, record.RelativePath()))
		} else {
			c := color.New(color.FgCyan)
			c.Printf("%v\n", record.RelativePath())
		}
	}

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
		Header:  MakeDumpHeader(path),
		Records: make([]DumpRecord, len(files)),
	}

	for i, file := range files {
		dFile.Records[i] = MakeDumpRecord(file)
	}

	//fmt.Println(dFile.Header)
	//fmt.Println(dFile.Records)

	b, err := msgpack.Marshal(dFile)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dumpFileName, b, 0644)
	if err != nil {
		return err
	}

	//var dFile2 DumpFile
	//dRecs2 := make([]DumpRecord, len(files))
	//msgpack.Unmarshal(b, &dFile2)
	//fmt.Println(dFile2)

	return nil
}

func loadRecords(dumpFileName string) (*DumpFile, error) {
	rawRecords, err := ioutil.ReadFile(dumpFileName)
	if err != nil {
		return nil, err
	}

	var dFile DumpFile

	err = msgpack.Unmarshal(rawRecords, &dFile)
	if err != nil {
		return nil, err
	}

	//fmt.Println(dFile)
	return &dFile, nil
}
