package main

import (
	"testing"
)

func TestGetDayOfWeekChar(t *testing.T) {
	tests := []struct {
		date   string
		format string
		want   string
	}{
		{"2023/08/02", "2006/01/02", "水"},
		{"2021/08/03", "2006/01/02", "火"},
		{"2021/08/07", "2006/01/02", "土"},
	}

	for _, tt := range tests {
		got, err := getDayOfWeekChar(tt.date, tt.format)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}
		if got != tt.want {
			t.Errorf("getDayOfWeekChar(%q, %q) = %q; want %q", tt.date, tt.format, got, tt.want)
		}
	}
}
