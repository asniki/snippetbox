package main

import (
	"asniki/snippetbox/internal/assert"
	"testing"
	"time"
)

func TestHumanDate2(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2025, 3, 17, 10, 15, 21, 0, time.UTC),
			want: "17 Mar 2025 at 10:15:21",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2025, 3, 17, 10, 15, 21, 0, time.FixedZone("CET", 1*60*60)),
			want: "17 Mar 2025 at 09:15:21",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)
			assert.Equal(t, hd, tt.want)
		})
	}
}
