package service

import (
	"strings"
	"testing"
)

func TestValidateExternalImageURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "valid", input: "https://images.example.com/map.webp", want: "https://images.example.com/map.webp"},
		{name: "trim and remove fragment", input: "  https://images.example.com/map.webp#preview  ", want: "https://images.example.com/map.webp"},
		{name: "empty", input: "   ", wantErr: true},
		{name: "http", input: "http://images.example.com/map.webp", wantErr: true},
		{name: "relative", input: "/map.webp", wantErr: true},
		{name: "credentials", input: "https://user:pass@images.example.com/map.webp", wantErr: true},
		{name: "localhost", input: "https://assets.localhost/map.webp", wantErr: true},
		{name: "loopback IPv4", input: "https://127.0.0.1/map.webp", wantErr: true},
		{name: "private IPv4", input: "https://192.168.1.10/map.webp", wantErr: true},
		{name: "loopback IPv6", input: "https://[::1]/map.webp", wantErr: true},
		{name: "too long", input: "https://images.example.com/" + strings.Repeat("a", 2049), wantErr: true},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got, err := validateExternalImageURL(test.input)
			if test.wantErr {
				if err == nil {
					t.Fatalf("validateExternalImageURL(%q) returned no error", test.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("validateExternalImageURL(%q) error = %v", test.input, err)
			}
			if got != test.want {
				t.Fatalf("validateExternalImageURL(%q) = %q, want %q", test.input, got, test.want)
			}
		})
	}
}
