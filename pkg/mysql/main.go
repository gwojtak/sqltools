package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

type ColumnDefinition struct {
	Name          string
	DataType      string
	Default       string
	Null          bool
	AutoIncrement bool
	Unique        bool
	Primary       bool
	Comment       string
}

type Table struct {
	Name             string
	DropBeforeCreate bool
}

type config struct {
	SqlDump *string
	Tables  []string
}

const start string = "CREATE TABLE `%s`"
const end string = "Table structure for table"

func using_strings(sql string, table string) string {
	var out string

	begin := fmt.Sprintf(start, table)
	in_table := false
	lines := strings.Split(sql, "\n")
	for _, line := range lines {
		if strings.Contains(line, begin) && in_table == false {
			in_table = true
			out += line + "\n"
		} else if strings.Contains(line, end) && in_table == true {
			in_table = false
			return out
		} else if in_table == true {
			out += line + "\n"
		}
	}
	return out
}

func readSQL(dumpConfig *config) string {
	var fh *os.File
	var err error
	var data string

	fmt.Println("Entering readSQL()")
	if *dumpConfig.SqlDump != "-" {
		fh, err = os.Open(*dumpConfig.SqlDump)
		if err != nil {
			panic(err)
		}
		defer fh.Close()
	} else {
		fh = os.Stdin
	}

	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		data += scanner.Text() + "\n"
	}

	return data
}

func parse_args() *config {
	dump_file := flag.String("f", "-", "sql dump file to search")
	flag.Parse()

	return &config{
		SqlDump: dump_file,
		Tables:  flag.Args(),
	}
}

func init() {}

func main() {
	config := parse_args()
	if len(config.Tables) < 1 {
		panic("You have to specify one or more tables")
	}
	sql := readSQL(config)
	for _, tbl := range config.Tables {
		fmt.Println(using_strings(sql, tbl))
	}
}
