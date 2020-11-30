package utils

import (
	"testing"

	"golang.org/x/net/context"
)

type validationValue struct {
	Value string `validate:"required,oneof=val"`
}

func TestMustValidate(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name    string
		args    args
		isError bool
	}{
		{"Should validate", args{&validationValue{"val"}}, false},
		{"Should fail", args{&validationValue{"v1"}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WrapPanic(func(ctx context.Context) {
				MustValidate(tt.args.val)
			})(context.Background()); tt.isError == (got == nil) {
				t.Errorf("MustValidate() error is %v, want error to be %v", got, tt.isError)
			}
		})
	}
}
