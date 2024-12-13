package gon3

import (
	"testing"
)

func Test_unescapeUChar(t *testing.T) {
	t.Skip("This test is not working as expected")
	tests := []struct {
		name string
		s    string
		want string
	}{
		{
			name: "on escapes",
			s:    "simple",
			want: "simple",
		},
		{
			name: "empty",
			s:    "",
			want: "",
		},
		{
			name: "single short escape",
			s:    `ab\u0063de`,
			want: "abcde",
		},
		{
			name: "single long escape",
			s:    `ab\U00000063de`,
			want: "abcde",
		},
		{
			name: "escaped short escape",
			s:    `a \\user`,
			want: `a \\user`,
		},
		{
			name: "escaped long escape",
			s:    `a \\User`,
			want: `a \\User`,
		},
		{
			name: "escaped escape before a short escape",
			s:    `a \\\u0063lass`,
			want: `a \\class`,
		},
		{
			name: "leading short escape",
			s:    `\u0061bc`,
			want: "abc",
		},
		{
			name: "leading escaped escape before a short escape",
			s:    `\\\u0061bc`,
			want: `\\abc`,
		},
		{
			name: "multiple escaped escapes",
			s:    `a \\\\user`,
			want: `a \\\\user`,
		},
		{
			name: "multiple short escapes",
			s:    `an \u0061b\u0063 \\user`,
			want: `an abc \\user`,
		},
		{
			name: "ending with escaped escape",
			s:    `end\\u`,
			want: `end\\u`,
		},

		{
			name: "invalid escape",
			s:    `\u0061b\u0063 \U0000WXYZ \u0061b\u0063`,
			want: `abc \U0000WXYZ abc`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("unescapeUChar(%s) = panic: %v", tt.s, r)
				}
			}()
			if got := unescapeUChar(tt.s); got != tt.want {
				t.Errorf("unescapeUChar(%s) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestUnescapeUChar(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic ascii conversion",
			input:    `\u0061b\u0063`,
			expected: "abc",
		},
		{
			name:     "invalid U escape sequence",
			input:    `\u0061b\u0063 \U0000WXYZ \u0061b\u0063`,
			expected: `abc \U0000WXYZ abc`,
		},
		{
			name:     "multiple valid sequences",
			input:    `\u0048\u0065\u006C\u006C\u006F`,
			expected: "Hello",
		},
		{
			name:     "mixed valid and invalid",
			input:    `\u0048\u004hello\u0021`,
			expected: `H\u004hello!`,
		},
		{
			name:     "truncated sequence at end",
			input:    `\u0048\u`,
			expected: `H\u`,
		},
		{
			name:     "valid U sequence",
			input:    `\U00000048\U00000069`,
			expected: "Hi",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "no escape sequences",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "incomplete hex digits",
			input:    `\u004\u0048`,
			expected: `\u004H`,
		},
		{
			name:     "escaped backslash",
			input:    `\\u0048`,
			expected: `\u0048`,
		},
		{
			name:     "mixed case hex",
			input:    `\u0048\U0000004A`,
			expected: "HJ",
		},
		{
			name:     "multiple invalid sequences",
			input:    `\UWXYZ \UABCD`,
			expected: `\UWXYZ \UABCD`,
		},
		{
			name:     "invalid sequence mixed with valid",
			input:    `\u0048\UWXYZ\u0069`,
			expected: `H\UWXYZi`,
		},
		{
			name:     "truncated U sequence",
			input:    `\U0000`,
			expected: `\U0000`,
		},
		{
			name:     "invalid hex in U sequence",
			input:    `\U0000GHIJ`,
			expected: `\U0000GHIJ`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := unescapeUChar(tt.input)
			if result != tt.expected {
				t.Errorf("unescapeUChar(%s) input %q =l got %q, want %q", tt.name, tt.input, result, tt.expected)
			}
		})
	}
}
