package random

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRandomStrings(t *testing.T) {
	tests := []struct {
		name string
		size int8
	}{
		{
			name: "size = 1",
			size: 1,
		},
		{
			name: "size = 5",
			size: 5,
		},
		{
			name: "size = 10",
			size: 10,
		},
		{
			name: "size = 30",
			size: 30,
		},
		{
			name: "size = 100",
			size: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str1 := NewRandomStrings(tt.size)
			str2 := NewRandomStrings(tt.size)

			assert.Len(t, str1, int(tt.size))
			assert.Len(t, str2, int(tt.size))

			assert.NotEqual(t, t, str1, str2)
		})
	}
}
