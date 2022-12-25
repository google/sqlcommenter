package core

import (
	"context"
	"testing"
)

func TestConvertMapToComment(t *testing.T) {
	for _, tc := range []struct {
		desc   string
		tagMap map[string]string
		want   string
	}{
		{
			desc: "nil tagMap",
			want: "",
		},
		{
			desc:   "no tags",
			tagMap: map[string]string{},
			want:   "",
		},
		{
			desc: "only one tag",
			tagMap: map[string]string{
				Route: "test-route",
			},
			want: "route='test-route'",
		},
		{
			desc: "only one tag with url encoding",
			tagMap: map[string]string{
				Route: "test/route",
			},
			want: "route='test%2Froute'",
		},
		{
			desc: "multiple tags",
			tagMap: map[string]string{
				Route:  "test/route",
				Action: "test-action",
				Driver: "sql-pg",
			},
			want: "action='test-action',db_driver='sql-pg',route='test%2Froute'",
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			if got, want := ConvertMapToComment(tc.tagMap), tc.want; got != want {
				t.Errorf("ConvertMapToComment(%+v) = %q, want = %q", tc.tagMap, got, want)
			}
		})
	}
}

type testRequestProvider struct {
	withRoute     string
	withAction    string
	withFramework string
}

func (p *testRequestProvider) Route() string     { return p.withRoute }
func (p *testRequestProvider) Action() string    { return p.withAction }
func (p *testRequestProvider) Framework() string { return p.withFramework }

func TestContextInject(t *testing.T) {
	tagsProvider := &testRequestProvider{
		withRoute:     "test-route",
		withAction:    "test-action",
		withFramework: "test-framework",
	}
	ctx := context.Background()
	gotCtx := ContextInject(ctx, tagsProvider)

	for _, tc := range []struct {
		desc string
		key  string
		want string
	}{
		{
			desc: "fetch action",
			key:  Action,
			want: "test-action",
		},
		{
			desc: "fetch route",
			key:  Route,
			want: "test-route",
		},
		{
			desc: "fetch framework",
			key:  Framework,
			want: "test-framework",
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			if got, want := gotCtx.Value(tc.key), tc.want; got != want {
				t.Errorf("ContextInject(ctx, tagsProvider) context.Value(%q) = %q, want = %q", tc.key, got, tc.want)
			}
		})
	}
}
