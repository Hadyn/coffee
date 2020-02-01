package jagex

import (
    "fmt"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestDecompressFileArchive(t *testing.T) {
    tests := []struct {
        name         string
        giveBytes    []byte
        wantBytes    []byte
        wantError    bool
        wantErrorMsg string
    }{
        {
            name:      "uncompressed",
            giveBytes: loadBytes(t, "archive.none.dat"),
            wantBytes: []byte{1},
        },
        {
            name:      "bzip",
            giveBytes: loadBytes(t, "archive.bzip.dat"),
            wantBytes: []byte("Hello World!"),
        },
        {
            name:      "gzip",
            giveBytes: loadBytes(t, "archive.gzip.dat"),
            wantBytes: []byte("Hello World!"),
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            defer func() {
                if r := recover(); r != nil {
                    t.Fatal("call panicked when did not want panic")
                }
            }()

            decompressed, err := DecompressFileArchive(tt.giveBytes)

            if tt.wantError {
                assert.EqualError(t, err, tt.wantErrorMsg)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.wantBytes, decompressed)
            }
        })
    }
}

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
            name: "uncompressed aligned",
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
            name: "uncompressed unaligned",
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
            name: "uncompressed aligned trailer",
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
        name         string
        give         []byte
        wantLength   int
        wantError    bool
        wantErrorMsg string
    }{
        {
            name:       "uncompressed",
            give:       []byte{ArchiveCompressionNone.AsByte(), 0, 0, 0, 1, 1},
            wantLength: 6,
            wantError:  false,
        },
        {
            name:       "bzip",
            give:       []byte{ArchiveCompressionBZIP.AsByte(), 0, 0, 0, 1, 0, 0, 0, 2, 1},
            wantLength: 10,
            wantError:  false,
        },
        {
            name:       "gzip",
            give:       []byte{ArchiveCompressionGZIP.AsByte(), 0, 0, 0, 1, 0, 0, 0, 2, 1},
            wantLength: 10,
            wantError:  false,
        },
    }

    for _, tt := range tests {
        t.Run(fmt.Sprintf("%s", tt.name), func(t *testing.T) {
            defer func() {
                if r := recover(); r != nil {
                    t.Fatal("call panicked when did not want panic")
                }
            }()

            length, err := FileArchiveLength(tt.give)

            if tt.wantError {
                assert.EqualError(t, err, tt.wantErrorMsg)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.wantLength, length)
            }
        })
    }
}
