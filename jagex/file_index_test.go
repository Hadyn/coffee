package jagex

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecodeFileIndex(t *testing.T) {
	tests := []struct {
		name          string
		give          []byte
		wantFileIndex *FileIndex
		wantError     bool
		wantErrorMsg  string
	}{
		{
			name: "simple",
			give: loadBytes(t, "index.simple.dat"),
			wantFileIndex: &FileIndex{
				Revision: 0,
				Groups: []*FileGroupEntry{
					{
						NameHash: 0,
						Checksum: 1,
						Revision: 1,
						Files: []*FileEntry{
							{NameHash: 0},
						},
					},
				},
			},
		},
		{
			name: "unsupported format",
			give: loadBytes(t, "index.bad-format.dat"),
			wantError: true,
			wantErrorMsg: "format is not supported: 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fi, err := DecodeFileIndex(tt.give)

			if tt.wantError {
				assert.EqualError(t, err, tt.wantErrorMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantFileIndex, fi)
			}
		})
	}
}
