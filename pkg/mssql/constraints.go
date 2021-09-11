package mssql

import (
	"database/sql"
	"fmt"
	"log"
)

const queryTemplate = `SELECT name, COL_NAME(parent_object_id, parent_column_id) AS column_name, type, definition FROM %s WHERE parent_object_id = OBJECT_ID('%s') AND name = '%s'`

type Constraint struct {
	Name       string
	Type       string
	Key        string
	Definition string
}

func NewConstraint(table string, constraint string) *Constraint {
	var stmt *sql.Stmt
	var err error
	var result *sql.Rows
	var c *Constraint

	defQuery := fmt.Sprintf(queryTemplate, "sys.default_constraints", table, constraint)
	stmt, err = DBConn.Prepare(defQuery)
	if err != nil {
		log.Println(err)
		return nil
	}

	result, err = stmt.Query()
	if err != nil {
		log.Println(err)
		return nil
	}

	result.Next()
	err = result.Scan(
		c.Name,
		c.Key,
		c.Type,
		c.Definition,
	)
	if err != nil {
		log.Println(err)
		return nil
	}

	return c
}

func (c *Constraint) SQLString(table string) string {
	switch c.Type {
	case "default":
		return fmt.Sprintf(`ALTER TABLE %s ADD CONSTRAINT %s DEFAULT '%s' FOR '%s'`, table, c.Name, c.Definition, c.Key)
	case "check":
		return fmt.Sprintf(`ALTER TABLE %s ADD CONSTRAINT %s CHECK %s`, table, c.Name, c.Definition)
	}
	return ""
}
