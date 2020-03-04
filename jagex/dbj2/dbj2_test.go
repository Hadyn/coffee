package dbj2

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
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
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.giveBytes), func(t *testing.T) {
			assert.Equal(t, tt.wantSum, Sum(tt.giveBytes))
		})
	}
}
