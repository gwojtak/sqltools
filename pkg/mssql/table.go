package mssql

import (
	"fmt"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
)

const (
	columnsQuery = `SELECT * FROM (SELECT c.name, TYPE_NAME(c.system_type_id) AS data_type, c.max_length, c.precision, c.scale, c.collation_name, c.is_nullable, c.is_rowguidcol, c.is_identity, isc.column_default
	FROM sys.columns AS c
		INNER JOIN %s.INFORMATION_SCHEMA.COLUMNS AS isc
		ON COL_NAME(OBJECT_ID('%s'), c.column_id) = isc.COLUMN_NAME
	WHERE c.object_id = OBJECT_ID('%s') AND isc.TABLE_NAME = '%s') AS i`
	identityQuery = `SELECT seed_value, increment_value
	FROM sys.identity_columns
	WHERE object_id = OBJECT_ID('%s') ) AS ic`
)

type Table struct {
	Name     string
	Identity string
	Columns  []Column
	Indexes  []Index
}

func NewTable(database string, table string) (*Table, error) {
	var identityColumn string

	info_schema := fmt.Sprintf(columnsQuery, database, table, table, table)

	stmt, err := DBConn.Prepare(info_schema)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	new_columns := []Column{}

	for result.Next() {
		tmp_col := Column{}
		tmp_datatype := ColumnDataType{}

		err := result.Scan(
			&tmp_col.Name,
			&tmp_datatype.TypeName,
			&tmp_datatype.MaxLength,
			&tmp_datatype.Precision,
			&tmp_datatype.Scale,
			&tmp_col.Collation,
			&tmp_col.Nullable,
			&tmp_col.IsRowGuidCol,
			&tmp_col.IsIdentity,
			&tmp_col.Default,
		)
		if err != nil {
			return nil, err
		}
		if tmp_col.IsIdentity.Bool() == true {
			fmt.Printf("\n\n%s %t\n\n", tmp_col.Name, tmp_col.IsIdentity.Bool())
			identityColumn = tmp_col.Name
			st, e := DBConn.Prepare(fmt.Sprintf("SELECT seed_value, increment_value FROM sys.identity_columns WHERE object_id = OBJECT_ID('%s')", table))
			if err != nil {
				return nil, err
			}

			r, e := st.Query()
			if e != nil {
				return nil, e
			}
			for r.Next() {
				err = r.Scan(&tmp_col.Seed, &tmp_col.Increment)
			}
		}
		tmp_col.Type = tmp_datatype
		new_columns = append(new_columns, tmp_col)
	}

	returnTable := Table{
		Name:     table,
		Columns:  new_columns,
		Identity: identityColumn,
	}

	returnTable.LoadIndexes()

	return &returnTable, nil
}

func (t *Table) LoadIndexes() error {
	var indexNames []string

	listIndexesQuery := fmt.Sprintf(`SELECT name FROM sys.indexes WHERE object_id = OBJECT_ID('%s')`, t.Name)

	stmt, err := DBConn.Prepare(listIndexesQuery)
	if err != nil {
		return err
	}
	result, err := stmt.Query()
	if err != nil {
		return err
	}

	for result.Next() {
		n := ""
		result.Scan(&n)
		indexNames = append(indexNames, n)
	}

	for _, n := range indexNames {
		i, err := NewIndex(t.Name, n)
		if err != nil {
			return err
		}
		t.Indexes = append(t.Indexes, *i)
	}
	return nil
}

func (t *Table) ListIndexes() ([]string, error) {
	var return_string []string

	stmt, err := DBConn.Prepare(fmt.Sprintf("SELECT name FROM sys.indexes WHERE is_hypothetical = 0 AND index_id != 0 AND object_id = OBJECT_ID('%s')", t.Name))
	if err != nil {
		return nil, err
	}
	result, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	for result.Next() {
		var x string
		result.Scan(&x)
		return_string = append(return_string, x)
	}
	return return_string, nil
}

func (t *Table) String() string {
	definition := fmt.Sprintf("CREATE TABLE %s (\n", t.Name)
	var columnDefinitions []string

	for _, col := range t.Columns {
		def := col.Definition()
		if def == "" {
			return ""
		}
		columnDefinitions = append(columnDefinitions, def)
	}

	// Hackity-hack: Fix this
	if len(t.Indexes[0].Name) > 1 {
		for _, idx := range t.Indexes {
			columnDefinitions = append(columnDefinitions, idx.String())
		}
	}
	definition = fmt.Sprintf("%s%s", definition, strings.Join(columnDefinitions, ",\n"))
	return fmt.Sprintf("%s\n);", definition)
}

func (t *Table) ColumnNames() []string {
	var cols []string
	for _, col := range t.Columns {
		cols = append(cols, col.Name)
	}

	return cols
}

func (t *Table) Dump() string {
	var returnString string
	query := `SELECT * FROM ` + t.Name
	columns := t.ColumnNames()
	ncols := len(columns)
	vals := make([]interface{}, ncols)

	var columnSpec []string
	for i := 0; i < ncols; i++ {
		vals[i] = new(string)
		columnSpec = append(columnSpec, columns[i])
	}
	result, err := DBConn.Query(query)
	if err != nil {
		return fmt.Sprintf("%s\n", err)
	}

	for result.Next() {
		row := make([]string, ncols)
		result.Scan(vals...)
		for i, v := range vals {
			if v != nil {
				row[i] = *(v.(*string))
			} else {
				row[i] = "NULL"
			}
		}
		returnString = fmt.Sprintf("%sINSERT INTO %s %s VALUES (%s);\n", returnString, t.Name, strings.Join(columnSpec, ", "), strings.Join(row, ", "))
	}
	return returnString
}

/*	// Indexes
	index := Index{}
	result.NextResultSet()
	for result.Next() {
		ret := result.Scan(
			&index.Name,
			&index.Description,
			&index.Keys,
		)

		if ret != nil {
			return nil, ret
		}
		new_table.Indexes = append(new_table.Indexes, index)
	}

	// Constraints
	constraint := Constraint{}
	result.NextResultSet()
	for result.Next() {
		ret := result.Scan(
			&constraint.Type,
			&constraint.Name,
			&constraint.DeleteAction,
			&constraint.UpdateAction,
			&constraint.StatusEnabled,
			&constraint.StatusForReplication,
			&constraint.Keys,
		)
		if ret != nil {
			log.Println("DEBUG DEBUG ERROR", ret)
			return nil, ret
		}
		new_table.Constraints = append(new_table.Constraints, constraint)
	}
*/
