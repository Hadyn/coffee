package jagex

import (
    "errors"
    "golang.org/x/crypto/xtea"
)

type ArchiveCompression uint8

const (
    FileArchiveCompressionNone ArchiveCompression = iota
    FileArchiveCompressionBZIP
    FileArchiveCompressionGZIP
)

func (ac ArchiveCompression) valid() error {
    if ac < FileArchiveCompressionNone || ac > FileArchiveCompressionGZIP {
        return errors.New("unrecognized compression")
    }
    return nil
}

func DecompressFileArchive(bs []byte) ([]byte, error) {
    return nil, nil
}

func DecryptFileArchive(bs []byte, key []byte) ([]byte, error) {
    cipher, err := xtea.NewCipher(key)
    if err != nil {
        return nil, err
    }

    cipher.Decrypt(nil, nil)
    return nil, nil
}