package main

import (
	"log"
	"os"

	"github.com/gwojtak/sqltools/pkg/mysql"
	"github.com/jessevdk/go-flags"
)

type Opts struct {
	File       string   `short:"f" long:"file" description:"name of the mysql dump file" default:"-"`
	Exclude    []string `short:"x" long:"exclude" description:"list of tables to exclude"`
	Positional struct {
		Tables []string `positional-arg-name:"table" required:"true"`
	} `positional-args:"true" required:"true"`
}

func main() {
	var opts Opts
	var parser *flags.Parser
	var err error

	parser = flags.NewParser(&opts, flags.HelpFlag|flags.PrintErrors)
	_, err = parser.Parse()
	if err != nil {
		parser.WriteHelp(os.Stderr)
		os.Exit(1)
	}

	filter, err := mysql.NewDumpFilter(opts.File, opts.Exclude, opts.Positional.Tables)
	if err != nil {
		log.Fatal(err)
	}
	filter.Stream()
}
