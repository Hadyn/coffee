package dbj2

import (
    "hash"
)

const (
    DBJ32Size = 4
    DBJ32BlockSize = 1
)

type dbj32 uint32

func New() hash.Hash32 {
    var s = dbj32(0)
    return &s
}

func Sum(b []byte) uint32 {
    h := New()
    _, _ = h.Write(b)
    return h.Sum32()
}

func (d *dbj32) Write(data []byte) (n int, err error) {
    v := *d
    for _, c := range data {
        v = (v << 5) - v + dbj32(c)
    }
    *d = v
    return len(data), nil
}

func (d *dbj32) Sum(in []byte) []byte {
    v := uint32(*d)
    return append(in, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}

func (d *dbj32) Reset() {
    *d = 0
}

func (d *dbj32) Sum32() uint32 {
    return uint32(*d)
}
