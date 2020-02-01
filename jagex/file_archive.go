package jagex

import (
    "bytes"
    "compress/gzip"
    "errors"
    "github.com/dsnet/compress/bzip2"
    "golang.org/x/crypto/xtea"
    "io"
)

const (
    minimumArchiveHeaderLength = 5
)

type ArchiveCompression byte

const (
    ArchiveCompressionNone ArchiveCompression = iota
    ArchiveCompressionBZIP
    ArchiveCompressionGZIP
)

func (c ArchiveCompression) Validate() error {
    if c < ArchiveCompressionNone || c > ArchiveCompressionGZIP {
        return errors.New("unrecognized compression")
    }
    return nil
}

func (c ArchiveCompression) HeaderLength() int {
    switch c {
    case ArchiveCompressionNone:
        return minimumArchiveHeaderLength
    case ArchiveCompressionBZIP, ArchiveCompressionGZIP:
        return minimumArchiveHeaderLength + 4
    default:
        panic("unhandled compression enumeration")
    }
}

func (c ArchiveCompression) AsByte() byte { return byte(c) }

func DecompressFileArchive(bs []byte) ([]byte, error) {
    rb := ReadBuffer(bs)

    var (
        c                = ArchiveCompression(rb.GetUint8())
        headerLength     = c.HeaderLength()
        compressedLength = rb.GetUint32AsInt()
    )

    switch c {
    case ArchiveCompressionNone:
        copied := make([]byte, compressedLength)
        copy(copied, bs[headerLength:headerLength+compressedLength])
        return copied, nil
    case ArchiveCompressionBZIP, ArchiveCompressionGZIP:
        decompressed := make([]byte, rb.GetUint32())

        var r io.Reader
        switch c {
        case ArchiveCompressionBZIP:
            var err error
            r, err = bzip2.NewReader(
                io.MultiReader(
                    bytes.NewReader([]byte("BZh9")),
                    bytes.NewReader(bs[headerLength:headerLength+compressedLength]),
                ),
                &bzip2.ReaderConfig{},
            )

            if err != nil {
                return nil, err
            }
        case ArchiveCompressionGZIP:
            var err error
            r, err = gzip.NewReader(
                bytes.NewReader(bs[headerLength: headerLength+compressedLength]),
            )

            if err != nil {
                return nil, err
            }
        }

        if _, err := io.ReadFull(r, decompressed); err != nil {
            return nil, err
        }

        return decompressed, nil
    default:
        panic("unhandled compression enumeration")
    }
}

func DecryptFileArchive(bs []byte, key []byte) ([]byte, error) {
    cipher, err := xtea.NewCipher(key)
    if err != nil {
        return nil, err
    }

    copied := make([]byte, len(bs))
    copy(copied, bs)

    archiveLength := FileArchiveLength(bs)
    for i := minimumArchiveHeaderLength; i < archiveLength-xtea.BlockSize; i += xtea.BlockSize {
        cipher.Decrypt(copied[i:], bs[i:i+xtea.BlockSize])
    }

    return copied, nil
}

func FileArchiveLength(bs []byte) int {
    rb := ReadBuffer(bs)
    return ArchiveCompression(rb.GetUint8()).HeaderLength() + rb.GetUint32AsInt()
}
