package mssql

import (
	"database/sql"
	"fmt"
	"net/url"

	version "github.com/gwojtak/sqltools/pkg/version"
)

var DBConn *sql.DB

type ConnectionConfig struct {
	Server   string
	Username string
	Password string
	Port     uint16
	Database string
}

func NewConnectionConfig(server string, port uint16, username string, password string, database string) *ConnectionConfig {
	return &ConnectionConfig{
		Server:   server,
		Port:     port,
		Username: username,
		Password: password,
		Database: database,
	}
}

func (c *ConnectionConfig) String() (string, error) {
	qs := url.Values{}
	qs.Add("database", c.Database)
	qs.Add("app name", fmt.Sprintf("wojo-mssql-%s", version.Version))

	u := &url.URL{
		Scheme:   "sqlserver",
		Host:     fmt.Sprintf("%s:%d", c.Server, c.Port),
		RawQuery: qs.Encode(),
	}

	if c.Username != "" && c.Password != "" {
		u.User = url.UserPassword(c.Username, c.Password)
	} else if c.Username != "" && c.Password == "" {
		u.User = url.User(c.Username)
	}

	return u.String(), nil
}
