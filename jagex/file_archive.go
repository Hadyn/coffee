package jagex

import (
    "bytes"
    "compress/gzip"
    "fmt"
    "github.com/dsnet/compress/bzip2"
    "golang.org/x/crypto/xtea"
    "io"
)

const (
    uncompressedArchiveHeaderLength = 5
    compressedArchiveHeaderLength   = 9
)

type ArchiveCompression byte

const (
    ArchiveCompressionNone ArchiveCompression = iota
    ArchiveCompressionBZIP
    ArchiveCompressionGZIP
)

func (c ArchiveCompression) Check() error {
    switch c {
    case ArchiveCompressionNone, ArchiveCompressionBZIP, ArchiveCompressionGZIP:
        return nil
    default:
        return fmt.Errorf("unrecognized compression: %d", c)
    }
}

func (c ArchiveCompression) HeaderLength() int {
    switch c {
    case ArchiveCompressionNone:
        return uncompressedArchiveHeaderLength
    case ArchiveCompressionBZIP, ArchiveCompressionGZIP:
        return compressedArchiveHeaderLength
    default:
        panic("unhandled compression enumeration")
    }
}

func (c ArchiveCompression) AsByte() byte { return byte(c) }

func DecompressFileArchive(bs []byte) ([]byte, error) {
    rb := ReadBuffer(bs)

    c := ArchiveCompression(rb.GetUint8())
    if err := c.Check(); err != nil {
        return nil, err
    }

    compressedLength := rb.GetUint32AsInt()

    switch c {
    case ArchiveCompressionNone:
        copied := make([]byte, compressedLength)
        copy(copied, rb.Get(compressedLength))
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
                    bytes.NewReader(rb.Get(compressedLength)),
                ),
                &bzip2.ReaderConfig{},
            )

            if err != nil {
                return nil, err
            }
        case ArchiveCompressionGZIP:
            var err error
            r, err = gzip.NewReader(
                bytes.NewReader(rb.Get(compressedLength)),
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

    archiveLength, err := FileArchiveLength(bs)
    if err != nil {
        return nil, err
    }

    for i := uncompressedArchiveHeaderLength; i < archiveLength-xtea.BlockSize+1; i += xtea.BlockSize {
        cipher.Decrypt(copied[i:], bs[i:i+xtea.BlockSize])
    }

    return copied, nil
}

func FileArchiveLength(bs []byte) (int, error) {
    rb := ReadBuffer(bs)

    c := ArchiveCompression(rb.GetUint8())
    if err := c.Check(); err != nil {
        return 0, err
    }

    return c.HeaderLength() + rb.GetUint32AsInt(), nil
}
