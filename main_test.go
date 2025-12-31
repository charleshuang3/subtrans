package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFromIndex(t *testing.T) {
	tests := []struct {
		input    string
		wantItem int
		wantLine int
		wantSeg  int
		wantErr  bool
	}{
		{"1,2,3", 1, 2, 3, false},
		{" 1 , 2 , 3 ", 1, 2, 3, false},
		{"1,2", 0, 0, 0, true},
		{"1,2,3,4", 0, 0, 0, true},
		{"a,2,3", 0, 0, 0, true},
		{"1,b,3", 0, 0, 0, true},
		{"1,2,c", 0, 0, 0, true},
	}

	for _, tt := range tests {
		gotItem, gotLine, gotSeg, err := parseFromIndex(tt.input)
		if tt.wantErr {
			assert.Error(t, err, "parseFromIndex(%q) should return error", tt.input)
		} else {
			assert.NoError(t, err, "parseFromIndex(%q) should not return error", tt.input)
			assert.Equal(t, tt.wantItem, gotItem, "item should match")
			assert.Equal(t, tt.wantLine, gotLine, "line should match")
			assert.Equal(t, tt.wantSeg, gotSeg, "seg should match")
		}
	}
}
