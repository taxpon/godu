package godu

type DumpHeader struct {
	Version string
	Path    string
}

type DumpRecord struct {
	Path string
	Size int64
}

type DumpFile struct {
	Header DumpHeader
	Records []DumpRecord
}


func NewDumpRecord(f *File) (*DumpRecord, error) {
	d := &DumpRecord{
		Path: f.RelativePath(),
		Size: f.Size(),
	}
	return d, nil
}

func HSize() {

}
