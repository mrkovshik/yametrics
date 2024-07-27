package util

import "testing"

func TestCutSeconds(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"1", args{"1s"}, 1, false},
		{"2", args{"2"}, 0, true},
		{"3", args{"3m"}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CutSeconds(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("CutSeconds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CutSeconds() got = %v, want %v", got, tt.want)
			}
		})
	}
}
