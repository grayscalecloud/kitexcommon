package tools

import (
	"reflect"
	"testing"
)

func TestFormatIds(t *testing.T) {
	type args struct {
		ids []int32
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "1", args: args{ids: []int32{1, 2, 3, 9, 11, 5, 10}}, want: []string{"1-3", "5", "9-11"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatIds(tt.args.ids); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FormatIds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseIds(t *testing.T) {
	type args struct {
		ids []string
	}
	tests := []struct {
		name string
		args args
		want []int32
	}{
		{name: "1", args: args{ids: []string{"1-3", "5", "9-11"}}, want: []int32{1, 2, 3, 5, 9, 10, 11}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseIds(tt.args.ids); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseIds() = %v, want %v", got, tt.want)
			}
		})
	}
}
