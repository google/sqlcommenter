package sqlcommenter

import (
	"strings"
	"testing"
)

func TestValues_String(t *testing.T) {
	tests := []struct {
		name string
		vs   Values
		want string
	}{
		{name: "nil", vs: nil, want: ""},
		{name: "empty", vs: Values{}, want: ""},
		{name: "empty cast", vs: Values(map[string]string{}), want: ""},
		{
			name: "drop empty key",
			vs:   Values(map[string]string{"": "val"}),
			want: "",
		},
		{
			name: "one",
			vs:   Values(map[string]string{"key": "val"}),
			want: "/*key='val'*/",
		},
		{
			name: "two",
			vs:   Values(map[string]string{"a": "1", "b": "2"}),
			want: "/*a='1',b='2'*/",
		},
		{
			name: "two reversed",
			vs:   Values(map[string]string{"b": "2", "a": "1"}), // technically, Go map iteration is random
			want: "/*a='1',b='2'*/",
		},
		{
			name: "name=DROP TABLE FOO",
			vs:   Values(map[string]string{"name": "DROP TABLE FOO"}),
			want: "/*name='DROP%20TABLE%20FOO'*/",
		},
		{
			name: `name''=DROP TABLE USERS'`,
			vs:   Values(map[string]string{"name''": `DROP TABLE USERS'`}),
			want: `/*name%27%27='DROP%20TABLE%20USERS%27'*/`,
		},
		{
			name: `exhibit`, // https://google.github.io/sqlcommenter/spec/#sql-commenter-exhibit
			vs: Values(map[string]string{
				"action":      `%2Fparam*d`,
				"controller":  `index`,
				"framework":   `spring`,
				"traceparent": `00-5bd66ef5095369c7b0d1f8f4bd33716a-c532cb4098ac3dd2-01`,
				"tracestate":  `congo%3Dt61rcWkgMzE%2Crojo%3D00f067aa0ba902b7`,
			}),
			want: "/*" + strings.Join([]string{
				"action='%252Fparam%2Ad'",
				"controller='index'",
				"framework='spring'",
				"traceparent='00-5bd66ef5095369c7b0d1f8f4bd33716a-c532cb4098ac3dd2-01'",
				"tracestate='congo%253Dt61rcWkgMzE%252Crojo%253D00f067aa0ba902b7'",
			}, ",") + "*/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.vs.String(); got != tt.want {
				t.Errorf("\nwant: %v\ngot:  %v", tt.want, got)
			}
		})
	}
}
