package utils

import (
	"context"
	"net/url"
	"testing"
)

func TestMustURLParse(t *testing.T) {
	type args struct {
		urlStr string
	}
	tests := []struct {
		name string
		args args
	}{
		{"Should parse URL without panic", args{"poroto://test.com/path?val=1"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var u *url.URL
			if err := WrapPanic(func(context.Context) {
				u = MustURLParse(tt.args.urlStr)
			})(context.TODO()); err != nil {
				t.Errorf("MustURLParse() panic with %v", err)
			}

			if u.String() != tt.args.urlStr {
				t.Errorf("MustURLParse() = %v, want %v", u.String(), tt.args.urlStr)
			}
		})
	}
}
