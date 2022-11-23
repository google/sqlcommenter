import "testing"

func TestNewHTTPRequestTags(t *testing.T) {
	rt := NewHTTPRequestTags("f", "r", "a")

	if rt.framework != "f" {
		t.Errorf("rt.framework - got: %s, want: %s", rt.framework, "f")
	}

	if rt.route != "r" {
		t.Errorf("rt.route - got: %s, want: %s", rt.route, "r")
	}

	if rt.action != "a" {
		t.Errorf("rt.action - got: %s, want: %s", rt.action, "r")
	}
}
