package tools

import "testing"

func TestConvImagePath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		baseUrl string
		want    string
	}{
		{
			name:    "相对路径",
			path:    "abcc.jpg",
			baseUrl: "https://example.com/",
			want:    "https://example.com/abcc.jpg",
		},
		{
			name:    "绝对路径(http)",
			path:    "http://other.com/abcc.jpg",
			baseUrl: "https://example.com/",
			want:    "http://other.com/abcc.jpg",
		},
		{
			name:    "绝对路径(https)",
			path:    "https://other.com/abcc.jpg",
			baseUrl: "https://example.com/",
			want:    "https://other.com/abcc.jpg",
		},
		{
			name:    "空路径",
			path:    "",
			baseUrl: "https://example.com/",
			want:    "https://example.com/",
		},
		{
			name:    "baseUrl不以/结尾",
			path:    "abcc.jpg",
			baseUrl: "https://example.com",
			want:    "https://example.com/abcc.jpg", // 现在修正为预期行为
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvImagePath(tt.path, tt.baseUrl); got != tt.want {
				t.Errorf("ConvImagePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
