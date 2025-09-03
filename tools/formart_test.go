package tools

import (
	"reflect"
	"testing"
)

func TestBoolFormat(t *testing.T) {
	type args struct {
		expr bool
		a    interface{}
		b    interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{name: "true", args: args{expr: true, a: "a", b: "b"}, want: "a"},
		{name: "true", args: args{expr: true, a: "a", b: "b"}, want: "a"},
		{name: "true", args: args{expr: true, a: 1, b: 2}, want: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BoolFormat(tt.args.expr, tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BoolFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}
