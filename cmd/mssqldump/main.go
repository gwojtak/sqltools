package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gwojtak/sqltools/pkg/mssql"

	flags "github.com/jessevdk/go-flags"
)

var opts struct {
	Database string `short:"d" long:"database" description:"name of the database to use" default:"default"`
	Server   string `short:"s" long:"server" description:"hostname or IP of SQL server to connect to" default:"localhost"`
	Port     uint16 `short:"p" long:"port" description:"port on remote server to connect to" default:"1433"`
	Username string `short:"u" long:"username" description:"name of the user to authenticate with"`
	Password string `short:"P" long:"password" description:"password to authenticate with"`
}

func main() {
	var conf *mssql.ConnectionConfig
	var t *mssql.Table
	var err error

	args, _ := flags.Parse(&opts)

	conf = mssql.NewConnectionConfig(opts.Server, opts.Port, opts.Username, opts.Password, opts.Database)
	connectionURI, err := conf.String()
	if err != nil {
		log.Fatal(err)
	}

	mssql.DBConn, err = sql.Open("sqlserver", connectionURI)
	if err != nil {
		log.Fatal(err)
	}
	defer mssql.DBConn.Close()

	t, err = mssql.NewTable(opts.Database, args[0])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(t.String())
	fmt.Println(t.Dump())
}
