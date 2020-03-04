package jagex

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMultipartFileCollapse(t *testing.T) {
	tests := []struct {
		name      string
		withFile  *MultipartFile
		wantBytes []byte
	}{
		{
			name:      "single chunk",
			withFile:  &MultipartFile{Parts: [][]byte{{0, 1, 2, 3, 4, 5}}},
			wantBytes: []byte{0, 1, 2, 3, 4, 5},
		},
		{
			name:      "multi chunk",
			withFile:  &MultipartFile{Parts: [][]byte{{0, 1, 2, 3, 4, 5}, {6, 7, 8, 9}}},
			wantBytes: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collapsed := tt.withFile.Collapse()
			assert.Equal(t, tt.wantBytes, collapsed)
		})
	}
}

func TestDecodeFileGroup(t *testing.T) {
	tests := []struct {
		name      string
		giveBytes []byte
		giveSize  int
		wantFiles []*MultipartFile
	}{
		{
			name:      "single file",
			giveBytes: []byte{1},
			giveSize:  1,
			wantFiles: []*MultipartFile{
				{Parts: [][]byte{{1}}},
			},
		},
		{
			name:      "multiple files single chunk",
			giveBytes: []byte{1, 2, 3, 0, 0, 0, 1, 0, 0, 0, 2, 1},
			giveSize:  2,
			wantFiles: []*MultipartFile{
				{Parts: [][]byte{{1}}},
				{Parts: [][]byte{{2, 3}}},
			},
		},
		{
			name: "multiple files multiple parts",
			giveBytes: []byte{
				1, 2, 3, 4, 5, 6,
				0, 0, 0, 1, 0, 0, 0, 2,
				0, 0, 0, 2, 0, 0, 0, 1,
				2,
			},
			giveSize: 2,
			wantFiles: []*MultipartFile{
				{Parts: [][]byte{{1}, {4, 5}}},
				{Parts: [][]byte{{2, 3}, {6}}},
			},
		},
		{
			name: "multiple files empty parts",
			giveBytes: []byte{
				1, 2, 3, 4, 5, 6,
				0, 0, 0, 1, 0, 0, 0, 2,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 2, 0, 0, 0, 1,
				3,
			},
			giveSize: 2,
			wantFiles: []*MultipartFile{
				{Parts: [][]byte{{1}, {}, {4, 5}}},
				{Parts: [][]byte{{2, 3}, {}, {6}}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files := DecodeFileGroup(tt.giveBytes, tt.giveSize)
			assert.Equal(t, tt.wantFiles, files)
		})
	}
}
