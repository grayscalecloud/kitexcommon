package tools

import (
	"reflect"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/grayscalecloud/kitexcommon/consts"
	"github.com/grayscalecloud/kitexcommon/model"
)

func TestPasetoAuth(t *testing.T) {
	type args struct {
		audience string
		pi       model.PasetoConfig
	}
	tests := []struct {
		name string
		args args
		want app.HandlerFunc
	}{
		{name: "test", args: args{audience: consts.User,
			pi: model.PasetoConfig{
				PubKey:   "4651eabf791920d32c8cf2e295d100030d369018305a64600583b63920f0ec4b",
				Implicit: "X25519",
			}}, want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PasetoAuth(tt.args.audience, tt.args.pi)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PasetoAuth() = %v, want %v", got, tt.want)
			}

		})
	}
}
func TestPasetoAuthNotForce(t *testing.T) {
	type args struct {
		audience string
		pi       model.PasetoConfig
	}
	tests := []struct {
		name string
		args args
		want app.HandlerFunc
	}{
		{name: "test", args: args{audience: consts.User,
			pi: model.PasetoConfig{
				PubKey:   "4651eabf791920d32c8cf2e295d100030d369018305a64600583b63920f0ec4b",
				Implicit: "X25519",
			}}, want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PasetoAuthNotForce(tt.args.audience, tt.args.pi)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PasetoAuth() = %v, want %v", got, tt.want)
			}

		})
	}
}
