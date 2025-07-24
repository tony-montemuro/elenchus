package main

import (
	"testing"
	"time"

	"github.com/tony-montemuro/elenchus/internal/assert"
)

func TestTimeAgo(t *testing.T) {
	tests := []struct {
		name      string
		timestamp time.Time
		expected  string
	}{
		{
			name:      "Seconds ago",
			timestamp: time.Now().Add(-10 * time.Second),
			expected:  "10 seconds ago",
		},
		{
			name:      "Minutes ago",
			timestamp: time.Now().Add(-20 * time.Minute),
			expected:  "20 minutes ago",
		},
		{
			name:      "Hours ago",
			timestamp: time.Now().Add(-4 * time.Hour),
			expected:  "4 hours ago",
		},
		{
			name:      "Days ago",
			timestamp: time.Now().Add(-2 * 24 * time.Hour),
			expected:  "2 days ago",
		},
		{
			name:      "Weeks ago",
			timestamp: time.Now().Add(-2 * 24 * 7 * time.Hour),
			expected:  "2 weeks ago",
		},
		{
			name:      "Months ago",
			timestamp: time.Now().Add(-7 * 24 * 30 * time.Hour),
			expected:  "7 months ago",
		},
		{
			name:      "Years ago",
			timestamp: time.Now().Add(-6 * 24 * 365 * time.Hour),
			expected:  "6 years ago",
		},
		{
			name:      "1 second",
			timestamp: time.Now().Add(-1 * time.Second),
			expected:  "1 second ago",
		},
		{
			name:      "1 minute",
			timestamp: time.Now().Add(-1 * time.Minute),
			expected:  "1 minute ago",
		},
		{
			name:      "1 hour",
			timestamp: time.Now().Add(-1 * time.Hour),
			expected:  "1 hour ago",
		},
		{
			name:      "1 day",
			timestamp: time.Now().Add(-1 * 24 * time.Hour),
			expected:  "1 day ago",
		},
		{
			name:      "1 week",
			timestamp: time.Now().Add(-1 * 24 * 7 * time.Hour),
			expected:  "1 week ago",
		},
		{
			name:      "1 month",
			timestamp: time.Now().Add(-1 * 24 * 28 * time.Hour),
			expected:  "1 month ago",
		},
		{
			name:      "1 year",
			timestamp: time.Now().Add(-1 * 24 * 365 * time.Hour),
			expected:  "1 year ago",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, timeAgo(tt.timestamp), tt.expected)
		})
	}
}
