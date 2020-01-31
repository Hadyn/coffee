package jagex

import (
    "errors"
    "golang.org/x/crypto/xtea"
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
    return nil, nil
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
