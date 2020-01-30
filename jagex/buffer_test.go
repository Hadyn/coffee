package jagex

import (
    "fmt"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestReadBufferGetUint8(t *testing.T) {
    tests := []struct {
        give       []byte
        wantValue  uint8
        wantLength int
        wantPanic  bool
    }{
        {
            give:       []byte{0},
            wantValue:  0,
            wantLength: 0,
            wantPanic:  false,
        },
        {
            give:       []byte{0, 0},
            wantValue:  0,
            wantLength: 1,
            wantPanic:  false,
        },
        {
            give:       []byte{1, 0},
            wantValue:  1,
            wantLength: 1,
            wantPanic:  false,
        },
        {
            give:      []byte{},
            wantPanic: true,
        },
    }

    for _, tt := range tests {
        t.Run(fmt.Sprintf("%v", tt.give), func(t *testing.T) {
            defer func() {
                if r := recover(); r != nil && !tt.wantPanic {
                    t.Fatal("call panicked when did not want panic")
                }
            }()

            rb := ReadBuffer(tt.give)
            assert.Equal(t, tt.wantValue, rb.GetUint8())
            assert.Equal(t, tt.wantLength, len(rb))
        })
    }
}

func TestReadBufferGetUint16(t *testing.T) {
    tests := []struct {
        give       []byte
        wantValue  uint16
        wantLength int
        wantPanic  bool
    }{
        {
            give:       []byte{0, 0},
            wantValue:  0,
            wantLength: 0,
        },
        {
            give:       []byte{0, 0, 0},
            wantValue:  0,
            wantLength: 1,
        },
        {
            give:       []byte{0, 1, 0},
            wantValue:  0x1,
            wantLength: 1,
        },
        {
            give:       []byte{1, 0, 0},
            wantValue:  0x100,
            wantLength: 1,
        },
        {
            give:       []byte{1, 1, 0},
            wantValue:  0x101,
            wantLength: 1,
        },
        {
            give:      []byte{},
            wantPanic: true,
        },
    }

    for _, tt := range tests {
        t.Run(fmt.Sprintf("%v", tt.give), func(t *testing.T) {
            defer func() {
                if r := recover(); r != nil && !tt.wantPanic {
                    t.Fatal("call panicked when did not want panic")
                }
            }()

            rb := ReadBuffer(tt.give)
            assert.Equal(t, tt.wantValue, rb.GetUint16())
            assert.Equal(t, tt.wantLength, len(rb))
        })
    }
}

func TestReadBufferGetUint24(t *testing.T) {
    tests := []struct {
        give       []byte
        wantValue  uint32
        wantLength int
        wantPanic  bool
    }{
        {
            give:       []byte{0, 0, 0},
            wantValue:  0,
            wantLength: 0,
        },
        {
            give:       []byte{0, 0, 0, 0},
            wantValue:  0,
            wantLength: 1,
        },
        {
            give:       []byte{0, 0, 1, 0},
            wantValue:  0x1,
            wantLength: 1,
        },
        {
            give:       []byte{0, 1, 0, 0},
            wantValue:  0x100,
            wantLength: 1,
        },
        {
            give:       []byte{1, 0, 0, 0},
            wantValue:  0x10000,
            wantLength: 1,
        },
        {
            give:       []byte{1, 1, 1, 0},
            wantValue:  0x10101,
            wantLength: 1,
        },
        {
            give:      []byte{},
            wantPanic: true,
        },
    }

    for _, tt := range tests {
        t.Run(fmt.Sprintf("%v", tt.give), func(t *testing.T) {
            defer func() {
                if r := recover(); r != nil && !tt.wantPanic {
                    t.Fatal("call panicked when did not want panic")
                }
            }()

            rb := ReadBuffer(tt.give)
            assert.Equal(t, tt.wantValue, rb.GetUint24())
            assert.Equal(t, tt.wantLength, len(rb))
        })
    }
}

func TestReadBufferGetUint32(t *testing.T) {
    tests := []struct {
        give       []byte
        wantValue  uint32
        wantLength int
        wantPanic  bool
    }{
        {
            give:       []byte{0, 0, 0, 0},
            wantValue:  0,
            wantLength: 0,
        },
        {
            give:       []byte{0, 0, 0, 0, 0},
            wantValue:  0,
            wantLength: 1,
        },
        {
            give:       []byte{0, 0, 0, 1, 0},
            wantValue:  0x1,
            wantLength: 1,
        },
        {
            give:       []byte{0, 0, 1, 0, 0},
            wantValue:  0x100,
            wantLength: 1,
        },
        {
            give:       []byte{0, 1, 0, 0, 0},
            wantValue:  0x10000,
            wantLength: 1,
        },
        {
            give:       []byte{1, 0, 0, 0, 0},
            wantValue:  0x1000000,
            wantLength: 1,
        },
        {
            give:       []byte{1, 1, 1, 1, 0},
            wantValue:  0x1010101,
            wantLength: 1,
        },
        {
            give:      []byte{},
            wantPanic: true,
        },
    }

    for _, tt := range tests {
        t.Run(fmt.Sprintf("%v", tt.give), func(t *testing.T) {
            defer func() {
                if r := recover(); r != nil && !tt.wantPanic {
                    t.Fatal("call panicked when did not want panic")
                }
            }()

            rb := ReadBuffer(tt.give)
            assert.Equal(t, tt.wantValue, rb.GetUint32())
            assert.Equal(t, tt.wantLength, len(rb))
        })
    }
}

func TestReadBufferGetUint64(t *testing.T) {
    tests := []struct {
        give       []byte
        wantValue  uint64
        wantLength int
        wantPanic  bool
    }{
        {
            give:       []byte{0, 0, 0, 0, 0, 0, 0, 0},
            wantValue:  0,
            wantLength: 0,
        },
        {
            give:       []byte{0, 0, 0, 0, 0, 0, 0, 0, 0},
            wantValue:  0,
            wantLength: 1,
        },
        {
            give:       []byte{0, 0, 0, 0, 0, 0, 0, 1, 0},
            wantValue:  0x1,
            wantLength: 1,
        },
        {
            give:       []byte{0, 0, 0, 0, 0, 0, 1, 0, 0},
            wantValue:  0x100,
            wantLength: 1,
        },
        {
            give:       []byte{0, 0, 0, 0, 0, 1, 0, 0, 0},
            wantValue:  0x10000,
            wantLength: 1,
        },
        {
            give:       []byte{0, 0, 0, 0, 1, 0, 0, 0, 0},
            wantValue:  0x1000000,
            wantLength: 1,
        },
        {
            give:       []byte{0, 0, 0, 1, 0, 0, 0, 0, 0},
            wantValue:  0x100000000,
            wantLength: 1,
        },
        {
            give:       []byte{0, 0, 1, 0, 0, 0, 0, 0, 0},
            wantValue:  0x10000000000,
            wantLength: 1,
        },
        {
            give:       []byte{0, 1, 0, 0, 0, 0, 0, 0, 0},
            wantValue:  0x1000000000000,
            wantLength: 1,
        },
        {
            give:       []byte{1, 0, 0, 0, 0, 0, 0, 0, 0},
            wantValue:  0x100000000000000,
            wantLength: 1,
        },
        {
            give:       []byte{1, 1, 1, 1, 1, 1, 1, 1, 0},
            wantValue:  0x101010101010101,
            wantLength: 1,
        },
        {
            give:      []byte{},
            wantPanic: true,
        },
    }

    for _, tt := range tests {
        t.Run(fmt.Sprintf("%v", tt.give), func(t *testing.T) {
            defer func() {
                if r := recover(); r != nil && !tt.wantPanic {
                    t.Fatal("call panicked when did not want panic")
                }
            }()

            rb := ReadBuffer(tt.give)
            assert.Equal(t, tt.wantValue, rb.GetUint64())
            assert.Equal(t, tt.wantLength, len(rb))
        })
    }
}

func TestReadBufferGetCString(t *testing.T) {
    tests := []struct {
        give       []byte
        wantValue  string
        wantLength int
        wantPanic  bool
    }{
        {
            give:       []byte("Hello World\000"),
            wantValue:  "Hello World",
            wantLength: 0,
        },
        {
            give:       []byte("\000"),
            wantValue:  "",
            wantLength: 0,
        },
        {
            give:       []byte("Hello World\000\000"),
            wantValue:  "Hello World",
            wantLength: 1,
        },
        {
            give:       []byte("Hello World"),
            wantPanic: true,
        },
    }

    for _, tt := range tests {
        t.Run(fmt.Sprintf("%s", tt.give), func(t *testing.T) {
            defer func() {
                if r := recover(); r != nil && !tt.wantPanic {
                    t.Fatal("call panicked when did not want panic")
                }
            }()

            rb := ReadBuffer(tt.give)
            assert.Equal(t, tt.wantValue, rb.GetCString())
            assert.Equal(t, tt.wantLength, len(rb))
        })
    }
}

func TestReadBufferGet(t *testing.T) {
    tests := []struct {
        giveBytes []byte
        giveN     int
        wantValue []byte
        wantLength int
        wantPanic  bool
    }{
        {
            giveBytes:  []byte{0, 0},
            giveN:      2,
            wantValue:  []byte{0, 0},
            wantLength: 0,
        },
        {
            giveBytes:  []byte{1, 1, 0},
            giveN:      2,
            wantValue:  []byte{1, 1},
            wantLength: 1,
        },
        {
            giveBytes:  []byte{1},
            giveN:      2,
            wantPanic: true,
        },
    }

    for _, tt := range tests {
        t.Run(fmt.Sprintf("%v", tt.giveBytes), func(t *testing.T) {
            defer func() {
                if r := recover(); r != nil && !tt.wantPanic {
                    t.Fatal("call panicked when did not want panic")
                }
            }()

            rb := ReadBuffer(tt.giveBytes)
            assert.Equal(t, tt.wantValue, rb.Get(tt.giveN))
            assert.Equal(t, tt.wantLength, len(rb))
        })
    }
}
