package jagex

import (
    "fmt"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestDbj32Sum(t *testing.T) {
    tests := []struct{
        giveBytes []byte
        wantSum   uint32
    } {
        {
            giveBytes: []byte{},
            wantSum:   0,
        },
    }

    for _, tt := range tests {
        t.Run(fmt.Sprintf("%v", tt.giveBytes), func(t *testing.T) {
            h := NewDBJ2()
            _, err := h.Write(tt.giveBytes)
            assert.NoError(t, err)
            assert.Equal(t, tt.wantSum, h.Sum32())
        })
    }
}
