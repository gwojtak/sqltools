package mssql

import (
	"fmt"
	"strings"
)

type Index struct {
	Name      string
	Columns   []string
	Clustered BetterBool
	Unique    BetterBool
}

func NewIndex(table string, index string) (*Index, error) {
	sql := `SELECT i.name AS index_name,
				COL_NAME(ic.object_id, ic.column_id) AS column_name,
				i.type_desc,
				i.is_unique,
				i.is_primary_key
			FROM sys.indexes AS i
			INNER JOIN sys.index_columns AS ic
				ON i.object_id = ic.object_id AND i.index_id = ic.index_id
			WHERE i.object_id = OBJECT_ID('%s') AND i.name = '%s'`
	var keys []string

	stmt, err := DBConn.Prepare(fmt.Sprintf(sql, table, index))
	if err != nil {
		return nil, err
	}

	result, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	var name string
	var column string
	var clustered BetterBool
	var unique BetterBool
	var pkey BetterBool

	for result.Next() {
		result.Scan(
			&name,
			&column,
			&clustered,
			&unique,
			&pkey,
		)
		keys = append(keys, column)
	}
	return &Index{
		Name:      name,
		Columns:   keys,
		Clustered: clustered,
		Unique:    unique,
	}, nil
}

func (i *Index) String() string {
	var clustered string
	var unique string

	if i.Unique.Bool() {
		unique = "UNIQUE"
	} else {
		unique = ""
	}
	if i.Clustered.Bool() {
		clustered = "CLUSTERED"
	} else {
		clustered = "NONCLUSTERED"
	}
	keys := strings.Join(i.Columns, ", ")

	idxRet := fmt.Sprintf("  INDEX %s %s %s (%s)", i.Name, unique, clustered, keys)

	return idxRet
}
