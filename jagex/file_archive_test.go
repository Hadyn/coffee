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
            wantBytes: []byte("Hello world!"),
        },
        {
            name:      "bzip",
            giveBytes: loadBytes(t, "archive.bzip.dat"),
            wantBytes: []byte("Hello world!"),
        },
        {
            name:      "gzip",
            giveBytes: loadBytes(t, "archive.gzip.dat"),
            wantBytes: []byte("Hello world!"),
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            defer func() {
                if r := recover(); r != nil {
                    t.Fatal("call panicked when did not want panic", r)
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
                ArchiveCompressionNone.AsByte(), 0x00, 0x00, 0x00, 0x08,
                0x67, 0xc5, 0x54, 0xa4, 0x7b, 0x6a, 0x75, 0xbf,
            },
            giveKey: []byte{
                0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
                0x08, 0x09, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05,
            },
            wantBytes: []byte{
                ArchiveCompressionNone.AsByte(), 0x00, 0x00, 0x00, 0x08,
                0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
            },
        },
        {
            name: "uncompressed unaligned",
            giveBytes: []byte{
                ArchiveCompressionNone.AsByte(),  0x00, 0x00, 0x00, 0x09,
                0x67, 0xc5, 0x54, 0xa4, 0x7b, 0x6a, 0x75, 0xbf, 0x02,
            },
            giveKey: []byte{
                0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
                0x08, 0x09, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05,
            },
            wantBytes: []byte{
                ArchiveCompressionNone.AsByte(), 0x00, 0x00, 0x00, 0x09,
                0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x02,
            },
        },
        {
            name: "uncompressed aligned trailer",
            giveBytes: []byte{
                ArchiveCompressionNone.AsByte(), 0x00, 0x00, 0x00, 0x08,
                0x67, 0xc5, 0x54, 0xa4, 0x7b, 0x6a, 0x75, 0xbf, 0x02,
            },
            giveKey: []byte{
                0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
                0x08, 0x09, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05,
            },
            wantBytes: []byte{
                ArchiveCompressionNone.AsByte(), 0x00, 0x00, 0x00, 0x08,
                0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x02,
            },
        },
        {
            name: "uncompressed unaligned trailer",
            giveBytes: []byte{
                ArchiveCompressionNone.AsByte(), 0x00, 0x00, 0x00, 0x09,
                0x67, 0xc5, 0x54, 0xa4, 0x7b, 0x6a, 0x75, 0xbf, 0x02, 0x03,
            },
            giveKey: []byte{
                0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
                0x08, 0x09, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05,
            },
            wantBytes: []byte{
                ArchiveCompressionNone.AsByte(), 0x00, 0x00, 0x00, 0x09,
                0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x02, 0x03,
            },
        },
        {
            name: "bzip aligned",
            giveBytes: []byte{
                ArchiveCompressionBZIP.AsByte(), 0x00, 0x00, 0x00, 0x04,
                0xc6, 0xf3, 0xe9, 0xc, 0x61, 0x20, 0x10, 0x9b,
            },
            giveKey: []byte{
                0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
                0x08, 0x09, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05,
            },
            wantBytes: []byte{
                ArchiveCompressionBZIP.AsByte(), 0x00, 0x00, 0x00, 0x04,
                0x00, 0x00, 0x00, 0x04, 0x01, 0x01, 0x01, 0x01,
            },
        },
        {
            name: "bzip unaligned",
            giveBytes: []byte{
                ArchiveCompressionBZIP.AsByte(),  0x00, 0x00, 0x00, 0x04,
                0xc6, 0xf3, 0xe9, 0xc, 0x61, 0x20, 0x10, 0x9b, 0x02,
            },
            giveKey: []byte{
                0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
                0x08, 0x09, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05,
            },
            wantBytes: []byte{
                ArchiveCompressionBZIP.AsByte(), 0x00, 0x00, 0x00, 0x04,
                0x00, 0x00, 0x00, 0x04, 0x01, 0x01, 0x01, 0x01, 0x02,
            },
        },
        {
            name: "bzip aligned trailer",
            giveBytes: []byte{
                ArchiveCompressionBZIP.AsByte(), 0x00, 0x00, 0x00, 0x04,
                0xc6, 0xf3, 0xe9, 0xc, 0x61, 0x20, 0x10, 0x9b, 0x02,
            },
            giveKey: []byte{
                0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
                0x08, 0x09, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05,
            },
            wantBytes: []byte{
                ArchiveCompressionBZIP.AsByte(), 0x00, 0x00, 0x00, 0x04,
                0x00, 0x00, 0x00, 0x04, 0x01, 0x01, 0x01, 0x01, 0x02,
            },
        },
        {
            name: "bzip unaligned trailer",
            giveBytes: []byte{
                ArchiveCompressionBZIP.AsByte(), 0x00, 0x00, 0x00, 0x04,
                0xc6, 0xf3, 0xe9, 0xc, 0x61, 0x20, 0x10, 0x9b, 0x02, 0x03,
            },
            giveKey: []byte{
                0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
                0x08, 0x09, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05,
            },
            wantBytes: []byte{
                ArchiveCompressionBZIP.AsByte(), 0x00, 0x00, 0x00, 0x04,
                0x00, 0x00, 0x00, 0x04, 0x01, 0x01, 0x01, 0x01, 0x02, 0x03,
            },
        },
        {
            name: "gzip aligned",
            giveBytes: []byte{
                ArchiveCompressionGZIP.AsByte(), 0x00, 0x00, 0x00, 0x04,
                0xc6, 0xf3, 0xe9, 0xc, 0x61, 0x20, 0x10, 0x9b,
            },
            giveKey: []byte{
                0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
                0x08, 0x09, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05,
            },
            wantBytes: []byte{
                ArchiveCompressionGZIP.AsByte(), 0x00, 0x00, 0x00, 0x04,
                0x00, 0x00, 0x00, 0x04, 0x01, 0x01, 0x01, 0x01,
            },
        },
        {
            name: "gzip unaligned",
            giveBytes: []byte{
                ArchiveCompressionGZIP.AsByte(),  0x00, 0x00, 0x00, 0x04,
                0xc6, 0xf3, 0xe9, 0xc, 0x61, 0x20, 0x10, 0x9b, 0x02,
            },
            giveKey: []byte{
                0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
                0x08, 0x09, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05,
            },
            wantBytes: []byte{
                ArchiveCompressionGZIP.AsByte(), 0x00, 0x00, 0x00, 0x04,
                0x00, 0x00, 0x00, 0x04, 0x01, 0x01, 0x01, 0x01, 0x02,
            },
        },
        {
            name: "gzip aligned trailer",
            giveBytes: []byte{
                ArchiveCompressionGZIP.AsByte(), 0x00, 0x00, 0x00, 0x04,
                0xc6, 0xf3, 0xe9, 0xc, 0x61, 0x20, 0x10, 0x9b, 0x02,
            },
            giveKey: []byte{
                0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
                0x08, 0x09, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05,
            },
            wantBytes: []byte{
                ArchiveCompressionGZIP.AsByte(), 0x00, 0x00, 0x00, 0x04,
                0x00, 0x00, 0x00, 0x04, 0x01, 0x01, 0x01, 0x01, 0x02,
            },
        },
        {
            name: "gzip unaligned trailer",
            giveBytes: []byte{
                ArchiveCompressionGZIP.AsByte(), 0x00, 0x00, 0x00, 0x04,
                0xc6, 0xf3, 0xe9, 0xc, 0x61, 0x20, 0x10, 0x9b, 0x02, 0x03,
            },
            giveKey: []byte{
                0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
                0x08, 0x09, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05,
            },
            wantBytes: []byte{
                ArchiveCompressionGZIP.AsByte(), 0x00, 0x00, 0x00, 0x04,
                0x00, 0x00, 0x00, 0x04, 0x01, 0x01, 0x01, 0x01, 0x02, 0x03,
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
