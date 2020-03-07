package dbj2

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDBJ2Sum(t *testing.T) {
	tests := []struct {
		giveBytes []byte
		wantSum   uint32
	}{
		{
			giveBytes: []byte{},
			wantSum:   0,
		},
		{
			giveBytes: []byte("Hello, World!"),
			wantSum:   1498789909,
		},
		{
			giveBytes: []byte("123456789"),
			wantSum:   2427588661,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.giveBytes), func(t *testing.T) {
			assert.Equal(t, tt.wantSum, Sum(tt.giveBytes))
		})
	}
}
