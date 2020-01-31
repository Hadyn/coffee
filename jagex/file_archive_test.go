package jagex

import (
    "fmt"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestDecryptFileArchive(t *testing.T) {
    tests := []struct {
        name         string
        giveBytes    []byte
        giveKey      []byte
        wantBytes    []byte
        wantError    bool
        wantErrorMsg string
    }{
        {
            name:      "uncompressed aligned",
            giveBytes: []byte{
                ArchiveCompressionNone.AsByte(), 0, 0, 0, 8, 0, 0, 0, 0, 0, 0, 0, 0,
            },
            giveKey: []byte{
                0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
            },
            wantBytes: []byte{
                ArchiveCompressionNone.AsByte(), 0, 0, 0, 8, 0, 0, 0, 0, 0, 0, 0, 0,
            },
        },
        {
            name:      "uncompressed unaligned",
            giveBytes: []byte{
                ArchiveCompressionNone.AsByte(), 0, 0, 0, 9, 0, 0, 0, 0, 0, 0, 0, 0, 0,
            },
            giveKey: []byte{
                0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
            },
            wantBytes: []byte{
                ArchiveCompressionNone.AsByte(), 0, 0, 0, 9, 0, 0, 0, 0, 0, 0, 0, 0, 0,
            },
        },
        {
            name:      "uncompressed aligned trailer",
            giveBytes: []byte{
                ArchiveCompressionNone.AsByte(), 0, 0, 0, 8, 0, 0, 0, 0, 0, 0, 0, 0, 1,
            },
            giveKey: []byte{
                0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
            },
            wantBytes: []byte{
                ArchiveCompressionNone.AsByte(), 0, 0, 0, 8, 0, 0, 0, 0, 0, 0, 0, 0, 1,
            },
        },
        {
            name: "uncompressed unaligned trailer",
            giveBytes: []byte{
                ArchiveCompressionNone.AsByte(), 0, 0, 0, 9, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
            },
            giveKey: []byte{
                0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
            },
            wantBytes: []byte{
                ArchiveCompressionNone.AsByte(), 0, 0, 0, 9, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            decrypted, err := DecryptFileArchive(tt.giveBytes, tt.giveKey)

            if tt.wantError {
                assert.EqualError(t, err, tt.wantErrorMsg)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.wantBytes, decrypted)
            }
        })
    }
}

func TestFileArchiveLength(t *testing.T) {
    tests := []struct {
        name       string
        give       []byte
        wantLength int
        wantPanic  bool
    }{
        {
            name:       "uncompressed",
            give:       []byte{ArchiveCompressionNone.AsByte(), 0, 0, 0, 1, 1},
            wantLength: 6,
            wantPanic:  false,
        },
        {
            name:       "bzip",
            give:       []byte{ArchiveCompressionBZIP.AsByte(), 0, 0, 0, 1, 0, 0, 0, 2, 1},
            wantLength: 10,
            wantPanic:  false,
        },
        {
            name:       "gzip",
            give:       []byte{ArchiveCompressionGZIP.AsByte(), 0, 0, 0, 1, 0, 0, 0, 2, 1},
            wantLength: 10,
            wantPanic:  false,
        },
        {
            name:      "insufficient bytes uncompressed",
            give:      []byte{ArchiveCompressionNone.AsByte(), 0, 0, 0,},
            wantPanic: true,
        },
        {
            name:      "insufficient bytes bzip",
            give:      []byte{ArchiveCompressionBZIP.AsByte(), 0, 0, 0,},
            wantPanic: true,
        },
        {
            name:      "insufficient bytes gzip",
            give:      []byte{ArchiveCompressionGZIP.AsByte(), 0, 0, 0,},
            wantPanic: true,
        },
        {
            name:      "unrecognized compression",
            give:      []byte{255, 0, 0, 0, 1, 1},
            wantPanic: true,
        },
    }

    for _, tt := range tests {
        t.Run(fmt.Sprintf("%s", tt.name), func(t *testing.T) {
            defer func() {
                if r := recover(); r != nil && !tt.wantPanic {
                    t.Fatal("call panicked when did not want panic")
                }
            }()

            assert.Equal(t, tt.wantLength, FileArchiveLength(tt.give))
        })
    }
}
