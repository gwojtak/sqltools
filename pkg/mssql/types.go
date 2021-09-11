package mssql

import (
	"fmt"
	"strconv"
	"strings"
)

type BetterBool struct {
	Value bool
	Valid bool // Valid is true if Bool is not nil
}

func (b *BetterBool) Bool() bool {
	return b.Value
}

func (b *BetterBool) Scan(value interface{}) error {
	if value == nil {
		b.Value, b.Valid = false, false
	}
	switch fmt.Sprintf("%T", value) {
	case "bool":
		ret, err := strconv.ParseBool(fmt.Sprintf("%t", value))
		if err != err {
			return err
		}
		b.Value = ret
		b.Valid = true
	case "string":
		t := fmt.Sprintf("%s", value)
		switch strings.ToLower(string(t)) {
		case "yes", "y", "true", "1":
			b.Value = true
			b.Valid = true
		case "no", "n", "false", "0":
			b.Value = false
			b.Valid = true
		}
	case "int8", "uint8", "int16", "uint16", "int32", "uint32", "int64", "uint64", "int", "uint", "uintptr", "float32", "float64", "complex64", "complex128":
		if value == 0 {
			b.Value = false
			b.Valid = true
		} else {
			b.Value = true
			b.Valid = true
		}
	}
	return nil
}
