package utl

import "testing"

func TestValidateAddress(t *testing.T) {
	tests := []struct {
		name string
		addr string
		want bool
	}{
		{"1", "127.0.0.1:8090", true},
		{"2", "localhost:8090", true},
		{"3", "ya.ru:8090", true},
		{"4", "localhost:809d0", false},
		{"5", "asdf", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateAddress(tt.addr); got != tt.want {
				t.Errorf("ValidateAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
