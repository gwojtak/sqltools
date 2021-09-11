package mssql

import (
	"database/sql"
	"fmt"
	"strings"
)

type DefaultValue struct {
	Value string
	Valid bool
}

func (dv *DefaultValue) Scan(value interface{}) error {
	var format string

	dtype := fmt.Sprintf("%T", value)
	if value == nil {
		dv.Value = ""
		dv.Valid = false
	}
	switch dtype {
	case "bool":
		format = "%t"
	case "string":
		format = "%s"
	default:
		format = "%d"
	}
	dv.Value = fmt.Sprintf(format, value)
	dv.Valid = true

	return nil
}

type ColumnDataType struct {
	TypeName string

	Precision sql.NullInt64
	Scale     sql.NullInt64

	MaxLength sql.NullInt64
}

func (cdt *ColumnDataType) String() string {
	lowerType := strings.ToLower(cdt.TypeName)

	switch lowerType {
	case "character", "char", "varchar":
		return fmt.Sprintf("%s(%d)", cdt.TypeName, cdt.MaxLength.Int64)
	case "float":
		return fmt.Sprintf("%s(%d)", cdt.TypeName, cdt.Precision.Int64)
	case "decimal", "dec", "numeric":
		return fmt.Sprintf("%s(%d,%d)", cdt.TypeName, cdt.Precision.Int64, cdt.Scale.Int64)
	case "clob", "character large object", "blob", "binary large object":
		return fmt.Sprintf("%s(%d)", cdt.TypeName, cdt.MaxLength.Int64)
	}

	return cdt.TypeName
}

type Column struct {
	Name         string
	Type         ColumnDataType
	Default      DefaultValue
	Nullable     BetterBool
	Collation    sql.NullString
	IsRowGuidCol BetterBool
	IsIdentity   BetterBool
	Seed         sql.NullInt64
	Increment    sql.NullInt64
}

func (c *Column) Definition() string {
	var stmt string

	stmt = fmt.Sprintf("  %s %s", c.Name, c.Type.String())
	if c.IsIdentity.Bool() == true {
		stmt = fmt.Sprintf("%s IDENTITY(%d, %d)", stmt, c.Seed.Int64, c.Increment.Int64)
	}

	if !c.Nullable.Bool() {
		stmt = fmt.Sprintf("%s NOT NULL", stmt)
	} else {
		stmt = fmt.Sprintf("%s NULL", stmt)
	}

	if c.IsRowGuidCol.Bool() {
		stmt = fmt.Sprintf("%s ROWGUIDCOL", stmt)
	}

	if c.Default.Valid == false {
		stmt = fmt.Sprintf("DEFAULT %s", c.Default.Value)
	}

	return stmt
}
