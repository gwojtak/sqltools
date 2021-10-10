package sqltools

import (
	"strconv"
	"strings"
)

/**
 * Set the version of the utilities globally.
 * TODO: Automate updating this in the build tools
 **/
type VersionString string

type VersionId struct {
	Major int
	Minor int
	Patch int
	Tag   string
}

const Version VersionString = "0.0.0-dev"

func GetVersion() VersionString {
	return Version
}

func GetVersionId() *VersionId {
	var fields []string
	var ver VersionId

	f := func(r rune) bool {
		return r == '.' || r == '-'
	}
	fields = strings.FieldsFunc(string(Version), f)
	M, _ := strconv.Atoi(fields[0])
	m, _ := strconv.Atoi(fields[1])
	p, _ := strconv.Atoi(fields[2])
	ver = VersionId{
		Major: M,
		Minor: m,
		Patch: p,
	}
	if len(fields) > 3 {
		ver.Tag = fields[3]
	}
	return &ver
}

func (v *VersionId) Equals(other *VersionId) bool {
	return v.Major == other.Major && v.Minor == other.Minor && v.Patch == other.Patch
}

func (v *VersionId) LessThan(other *VersionId) bool {
	if v.Major < other.Major {
		return true
	} else if v.Major == other.Major && v.Minor < other.Minor {
		return true
	} else if v.Major == other.Major && v.Minor == other.Minor && v.Patch < other.Patch {
		return true
	}
	return false
}

func (v *VersionId) LessThanOrEqual(other *VersionId) bool {
	if v.Major < other.Major {
		return true
	} else if v.Major == other.Major && v.Minor < other.Minor {
		return true
	} else if v.Major == other.Major && v.Minor == other.Minor && v.Patch < other.Patch {
		return true
	} else if v.Equals(other) {
		return true
	}

	return false
}

func (v *VersionId) GreterThan(other *VersionId) bool {
	if v.Major > other.Major {
		return true
	} else if v.Major == other.Major && v.Minor > other.Minor {
		return true
	} else if v.Major == other.Major && v.Minor == other.Minor && v.Patch > other.Patch {
		return true
	}

	return false
}

func (v *VersionId) GreterThanOrEqual(other *VersionId) bool {
	if v.Major > other.Major {
		return true
	} else if v.Major == other.Major && v.Minor > other.Minor {
		return true
	} else if v.Major == other.Major && v.Minor == other.Minor && v.Patch > other.Patch {
		return true
	} else if v.Equals(other) {
		return true
	}

	return false
}
