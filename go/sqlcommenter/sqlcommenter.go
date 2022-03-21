package sqlcommenter

import (
	"net/url"
	"sort"
	"strings"
)

// Values maps a string key to a value for that key to attach to a SQL query
// in a comment. Implements the SQL Commenter spec:
// https://google.github.io/sqlcommenter
type Values map[string]string

// String returns the string representing all values according to the SQL
// Commenter spec.
func (vs Values) String() string {
	if len(vs) == 0 {
		return ""
	}

	pairs := make([]string, 0, len(vs))
	for k, v := range vs {
		if k == "" {
			continue
		}
		pairs = append(pairs, serializeKey(k)+"="+serializeValue(v))
	}

	if len(pairs) == 0 {
		return "" // we might have dropped only empty keys
	}

	// Spec requires sorted key-value pairs after running the serialization
	// algorithm.
	sort.Strings(pairs)

	return "/*" + strings.Join(pairs, ",") + "*/"
}

// https://google.github.io/sqlcommenter/spec/#key-serialization-algorithm
func serializeKey(s string) string {
	esc := urlEncode(s)
	return escapeMeta(esc)
}

// https://google.github.io/sqlcommenter/spec/#value-serialization-algorithm
func serializeValue(s string) string {
	esc := urlEncode(s)
	return `'` + escapeMeta(esc) + `'`
}

func urlEncode(s string) string {
	esc := url.QueryEscape(s)
	// Go encodes spaces as "+"; use more standard %20.
	return strings.Replace(esc, "+", "%20", -1)
}

func escapeMeta(s string) string {
	return strings.Replace(s, `'`, `\'`, -1)
}
