package mssql

func Compact(in []string) []string {
	var t = []string{}
	for _, v := range in {
		if v != "" {
			t = append(t, v)
		}
	}
	return t
}

func Contains(needle string, haystack []string) bool {
	for _, v := range haystack {
		if needle == v {
			return true
		}
	}
	return false
}
