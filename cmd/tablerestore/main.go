package main

import (
	"log"
	"os"
	"path"

	"github.com/gwojtak/sqltools/pkg/filtereddump"
	_ "github.com/gwojtak/sqltools/pkg/mssql"
	"github.com/jessevdk/go-flags"
)

type Opts struct {
	File       string `short:"f" long:"file" description:"name of the mysql dump file" default:"-"`
	Exclude    bool   `short:"v" description:"invert the match, ie. dump all tables except those specified"`
	Positional struct {
		Tables []string `positional-arg-name:"table" required:"true"`
	} `positional-args:"true" required:"true"`
}

func main() {
	var opts Opts
	var parser *flags.Parser
	var err error
	var dbtype string
	var bufferlen int

	switch path.Base(os.Args[0]) {
	case "mytablerestore":
		dbtype = "mysql"
		break
	case "mstablerestore":
		dbtype = "mssql"
		break
	}
	bufferlen = filtereddump.DefaultMaxBufferLen
	parser = flags.NewParser(&opts, flags.HelpFlag|flags.PrintErrors)
	_, err = parser.Parse()
	if err != nil {
		parser.WriteHelp(os.Stderr)
		os.Exit(1)
	}

	filter, err := filtereddump.NewDumpFilter(opts.File, opts.Exclude, opts.Positional.Tables, bufferlen, dbtype)
	if err != nil {
		log.Fatal(err)
	}
	filter.Stream()
}
