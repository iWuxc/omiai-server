package match

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTime(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "RFC3339",
			input:   "2023-10-27T10:00:00Z",
			wantErr: false,
		},
		{
			name:    "YYYY-MM-DD HH:mm:ss",
			input:   "2023-10-27 10:00:00",
			wantErr: false,
		},
		{
			name:    "YYYY-MM-DD",
			input:   "2023-10-27",
			wantErr: false,
		},
		{
			name:    "Invalid",
			input:   "invalid-date",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTime(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.False(t, got.IsZero())
			}
		})
	}
}
