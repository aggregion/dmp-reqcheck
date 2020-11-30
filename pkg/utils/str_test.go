package utils

import "testing"

func TestCoalesce(t *testing.T) {
	type args struct {
		strArgs []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Should get empty", args{[]string{""}}, ""},
		{"Should get first string", args{[]string{"first", ""}}, "first"},
		{"Should get third non empty string", args{[]string{"", "", "third", "fourth"}}, "third"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Coalesce(tt.args.strArgs...); got != tt.want {
				t.Errorf("Coalesce() = %v, want %v", got, tt.want)
			}
		})
	}
}
