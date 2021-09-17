package mysql

import (
	"bufio"
	"os"
	"strings"
)

const (
	DumpMarkerStart string = "CREATE TABLE "
	DumpMarkerEnd   string = "-- Table structure for table"
	MaxBufferLen    int    = 4 * 1024 * 1024
)

type DumpFilter struct {
	Input  *os.File
	Tables []string
}

func NewDumpFilter(filename string, exclude []string, tables []string) (*DumpFilter, error) {
	var t []string
	var i *os.File
	var err error

	if filename == "-" {
		i = os.Stdin
	} else {
		i, err = os.Open(filename)
		if err != nil {
			return nil, err
		}
	}

	for _, table := range tables {
		keepTable := true
		for _, ex := range exclude {
			if ex == table {
				keepTable = false
			}
		}
		if keepTable {
			t = append(t, table)
		}
	}

	return &DumpFilter{
		Input:  i,
		Tables: t,
	}, nil
}

func (d *DumpFilter) Stream() {
	var buffer string
	var scanner *bufio.Scanner
	var inTable bool
	var reader *bufio.Reader
	var buf []byte

	reader = bufio.NewReaderSize(d.Input, MaxBufferLen)
	scanner = bufio.NewScanner(reader)
	buf = make([]byte, MaxBufferLen)
	scanner.Buffer(buf, MaxBufferLen)
	scanner.Split(bufio.ScanLines)
	inTable = false

	for scanner.Scan() {
		buffer = scanner.Text()
		if strings.HasPrefix(buffer, DumpMarkerStart) {
			t := strings.Split(buffer, "`")[1]
			if IsIn(t, d.Tables) {
				inTable = true
			}
		}
		if strings.HasPrefix(buffer, DumpMarkerEnd) || strings.HasPrefix(buffer, "-- Final view structure") {
			inTable = false
		}
		if inTable {
			os.Stdout.Write([]byte(buffer + "\n"))
		}
	}
}

func RemoveElement(slice []string, value string) []string {
	var ret []string
	var found bool

	ret = make([]string, len(slice)-1)
	found = false
	for idx, val := range slice {
		if idx == len(slice)-1 && val != value {
			return slice
		}
		if val == value {
			found = true
		} else {
			if found {
				ret[idx-1] = val
			} else {
				ret[idx] = val
			}
		}
	}
	return ret
}

func IsIn(needle string, haystack []string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}
