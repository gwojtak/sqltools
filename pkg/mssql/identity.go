package mssql

import (
	"fmt"
)

type Identity struct {
	Column    string
	Seed      int
	Increment int
}

func NewIdentity(table string) (*Identity, error) {
	ident_query := fmt.Sprintf(`SELECT name, seed_value, increment_value FROM sys.identity_columns WHERE is_identity = 1 AND object_id = OBJECT_ID('%s')`, table)

	stmt, err := DBConn.Prepare(ident_query)
	if err != nil {
		return nil, err
	}

	result, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	var identity Identity
	result.Next()
	err = result.Scan(&identity.Column, &identity.Seed, &identity.Increment)
	if err != nil {
		return nil, err
	}
	return &identity, nil
}
