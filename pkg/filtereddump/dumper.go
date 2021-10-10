package filtereddump

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

const DefaultMaxBufferLen int = 4 * 1024 * 1024

var BufferLen int = DefaultMaxBufferLen

type TableMarker struct {
	DumpStart string
	DumpEnd   string
	ViewSkip  string
}

type FilteredDump struct {
	Input        *os.File
	Invert       bool
	Tables       []string
	MaxBufferLen int
	TableMarker
}

func NewDumpFilter(fd string, invert bool, tables []string, bufsize int, dbtype string) (*FilteredDump, error) {
	var in *os.File
	var err error

	if fd == "-" {
		in = os.Stdin
	} else {
		in, err = os.Open(fd)
		if err != nil {
			return nil, err
		}
	}

	if bufsize < 65536 {
		bufsize = 65536
	}

	d := FilteredDump{
		Input:        in,
		Invert:       invert,
		Tables:       tables,
		MaxBufferLen: bufsize,
	}
	d.SetDumpType(dbtype)
	return &d, nil
}

func (d *FilteredDump) Stream() {
	var buffer string
	var scanner *bufio.Scanner
	var reader *bufio.Reader
	var buf []byte
	var doWrite bool

	reader = bufio.NewReaderSize(d.Input, d.MaxBufferLen)
	scanner = bufio.NewScanner(reader)
	buf = make([]byte, d.MaxBufferLen)
	scanner.Buffer(buf, d.MaxBufferLen)
	scanner.Split(bufio.ScanLines)
	doWrite = false

	for scanner.Scan() {
		buffer = scanner.Text()
		if strings.HasPrefix(buffer, d.TableMarker.DumpStart) {
			t := strings.Split(buffer, "`")[1]
			flagged := IsIn(t, d.Tables)
			if flagged {
				if !d.Invert {
					doWrite = true
				} else {
					doWrite = false
				}
				d.RemoveTable(t)
			}
		}
		if strings.HasPrefix(buffer, d.TableMarker.DumpEnd) || strings.HasPrefix(buffer, d.TableMarker.ViewSkip) {
			if len(d.Tables) == 0 {
				return
			}
			doWrite = false
		}
		if doWrite {
			if !strings.HasPrefix(buffer, "--") {
				os.Stdout.Write([]byte(buffer + "\n"))
			}
		}
	}
}

func (f *FilteredDump) RemoveTable(value string) {
	var sL int
	var found int

	sL = len(f.Tables)
	if sL == 0 {
		return
	}
	if sL == 1 {
		f.Tables[0] = ""
		return
	}
	for idx, val := range f.Tables {
		if val == value {
			found = idx
		}
	}
	switch found {
	case 0:
		// The first element
		f.Tables = f.Tables[1:]
		return
	case (sL - 2):
		// The second to last element
		f.Tables = append(f.Tables[:(sL-2)], f.Tables[(sL-1)])
		return
	case (sL - 1):
		// The last element
		f.Tables = f.Tables[0:(sL - 1)]
		return
	default:
		// Somewhere in the middle
		f.Tables = append(f.Tables[:found], f.Tables[(found+1):]...)
		return
	}
}

func (d *FilteredDump) SetDumpType(dt string) error {
	switch dt {
	case "mysql":
		d.DumpStart = "CREATE TABLE "
		d.DumpEnd = "-- Table structure for table"
		d.ViewSkip = "-- Final view structure"
		return nil
	case "mssql":
		d.DumpStart = "-- Dumping table structure and data of table "
		d.DumpEnd = "-- End table dump"
		return nil
	default:
		break
	}
	return errors.New("unknown db type or none specified")
}

func IsIn(needle string, haystack []string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}
