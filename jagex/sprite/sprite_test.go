package sprite

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestDecodeSheet(t *testing.T) {
    tests := []struct {
        name      string
        giveBytes []byte
        wantSheet *Sheet
    }{
        {
            name:      "horizontally-encoded",
            giveBytes: []byte{
                0x00,
                0x00, 0x01, 0x01, 0x00, 0x00, 0x01,
                0x00, 0x00, 0x01,
                0x00, 0x04, 0x00, 0x04,
                0x01,
                0x00, 0x00,
                0x00, 0x01,
                0x00, 0x02,
                0x00, 0x03,
                0x00, 0x01,
            },
            wantSheet: &Sheet{
                Width:        4,
                Height:       4,
                OffsetX:      []int{0},
                OffsetY:      []int{1},
                PackedWidth:  []int{2},
                PackedHeight: []int{3},
                Colors:       []uint32{0, 1},
                Index:        [][]uint8{
                    {
                        0, 1,
                        1, 0,
                        0, 1,
                    },
                },
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            sheet := DecodeSheet(tt.giveBytes)
            assert.Equal(t, tt.wantSheet, sheet)
        })
    }
}
