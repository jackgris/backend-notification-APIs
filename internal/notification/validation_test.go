package notification

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateCategories(t *testing.T) {
	tests := []struct {
		name     string
		category string
		expected bool
	}{
		{
			name:     "Valid category Sports",
			category: "Sports",
			expected: true,
		},
		{
			name:     "Valid category Finance",
			category: "Finance",
			expected: true,
		},
		{
			name:     "Valid category Films",
			category: "Films",
			expected: true,
		},
		{
			name:     "Invalid category Music",
			category: "Music",
			expected: false,
		},
		{
			name:     "Empty category",
			category: "",
			expected: false,
		},
		{
			name:     "Case-sensitive test for sports",
			category: "sports",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validateCategories(tt.category)
			assert.Equal(t, tt.expected, result)
		})
	}
}
