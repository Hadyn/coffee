package jagex

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestCacheReaderGet(t *testing.T) {
	tests := []struct {
		name         string
		withLookup   io.ReadSeeker
		withBlocks   io.ReadSeeker
		withFileType int
		giveID       int
		wantBytes    []byte
		wantError    bool
		wantErrorMsg string
	}{
		{
			name:         "short",
			withLookup:   bytes.NewReader(loadBytes(t, "cache/lookup.short.dat")),
			withBlocks:   bytes.NewReader(loadBytes(t, "cache/blocks.short.dat")),
			withFileType: 1,
			giveID:       1,
			wantBytes:    loadBytes(t, "cache/short.dat"),
		},
		{
			name:         "long",
			withLookup:   bytes.NewReader(loadBytes(t, "cache/lookup.long.dat")),
			withBlocks:   bytes.NewReader(loadBytes(t, "cache/blocks.long.dat")),
			withFileType: 1,
			giveID:       1,
			wantBytes:    loadBytes(t, "cache/long.dat"),
		},
		{
			name:         "skip block",
			withLookup:   bytes.NewReader(loadBytes(t, "cache/lookup.skip-over.dat")),
			withBlocks:   bytes.NewReader(loadBytes(t, "cache/blocks.skip-over.dat")),
			withFileType: 1,
			giveID:       1,
			wantBytes:    loadBytes(t, "cache/skip-over.dat"),
		},
		{
			name:         "identifier mismatch",
			withLookup:   bytes.NewReader(loadBytes(t, "cache/lookup.bad-identifier.dat")),
			withBlocks:   bytes.NewReader(loadBytes(t, "cache/blocks.bad-identifier.dat")),
			withFileType: 1,
			giveID:       1,
			wantError:    true,
			wantErrorMsg: "file identifier mismatch; expected: 1, found: 2",
		},
		{
			name:         "chunk mismatch",
			withLookup:   bytes.NewReader(loadBytes(t, "cache/lookup.bad-chunk.dat")),
			withBlocks:   bytes.NewReader(loadBytes(t, "cache/blocks.bad-chunk.dat")),
			withFileType: 1,
			giveID:       1,
			wantError:    true,
			wantErrorMsg: "file chunk mismatch; expected: 1, found: 2",
		},
		{
			name:         "file type mismatch",
			withLookup:   bytes.NewReader(loadBytes(t, "cache/lookup.bad-file-type.dat")),
			withBlocks:   bytes.NewReader(loadBytes(t, "cache/blocks.bad-file-type.dat")),
			withFileType: 1,
			giveID:       1,
			wantError:    true,
			wantErrorMsg: "file type mismatch; expected: 1, found: 2",
		},
		{
			name:         "premature eof",
			withLookup:   bytes.NewReader(loadBytes(t, "cache/lookup.bad-next-block.dat")),
			withBlocks:   bytes.NewReader(loadBytes(t, "cache/blocks.bad-next-block.dat")),
			withFileType: 1,
			giveID:       1,
			wantError:    true,
			wantErrorMsg: "reached unexpected EOF; write-offset: 512, file-length: 513",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := NewCacheReader(tt.withLookup, tt.withBlocks, tt.withFileType)
			bs, err := fc.Read(tt.giveID)

			if tt.wantError {
				assert.EqualError(t, err, tt.wantErrorMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantBytes, bs)
			}
		})
	}
}
