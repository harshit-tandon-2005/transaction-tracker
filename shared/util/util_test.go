package util

import (
	"testing"
)

func TestFormatUnixTimestampString(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      string
		expectErr bool
	}{
		{
			name:      "Valid Timestamp",
			input:     "1710298091",
			want:      "2024-03-13 02:48:11",
			expectErr: false,
		},
		{
			name:      "Zero Timestamp",
			input:     "0",
			want:      "00-00-0000 00:00:00",
			expectErr: false,
		},
		{
			name:      "Invalid Input - Non-numeric",
			input:     "not-a-timestamp",
			want:      "",
			expectErr: true,
		},
		{
			name:      "Empty Input",
			input:     "",
			want:      "",
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := FormatUnixTimestampString(tc.input)

			if tc.expectErr {
				if err == nil {
					t.Errorf("FormatUnixTimestampString(%q) expected an error, but got nil", tc.input)
				}
			} else {
				if err != nil {
					t.Errorf("FormatUnixTimestampString(%q) unexpected error: %v", tc.input, err)
				}
			}

			if !tc.expectErr && got != tc.want {
				t.Errorf("FormatUnixTimestampString(%q) = %q; want %q", tc.input, got, tc.want)
			}
		})
	}
}
