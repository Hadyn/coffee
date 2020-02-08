package jagex

type ReadBuffer []byte

func (rb *ReadBuffer) GetUint8() (val uint8) {
    _ = (*rb)[0]
    val = (*rb)[0]
    *rb = (*rb)[1:]
    return
}

func (rb *ReadBuffer) GetUint8AsInt() (val int) {
    return int(rb.GetUint8())
}

func (rb *ReadBuffer) GetUint16() (val uint16) {
    _ = (*rb)[1]
    val = uint16((*rb)[0])<<8 | uint16((*rb)[1])
    *rb = (*rb)[2:]
    return
}

func (rb *ReadBuffer) GetUint16AsInt() (val int) {
	return int(rb.GetUint16())
}

func (rb *ReadBuffer) GetUint24() (val uint32) {
    _ = (*rb)[2]
    val = uint32((*rb)[0])<<16 | uint32((*rb)[1])<<8 | uint32((*rb)[2])
    *rb = (*rb)[3:]
    return
}

func (rb *ReadBuffer) GetUint24AsInt() (val int) {
	return int(rb.GetUint24())
}

func (rb *ReadBuffer) GetUint32() (val uint32) {
    _ = (*rb)[3]
    val = uint32((*rb)[0])<<24 | uint32((*rb)[1])<<16 | uint32((*rb)[2])<<8 | uint32((*rb)[3])
    *rb = (*rb)[4:]
    return
}

func (rb *ReadBuffer) GetUint32AsInt() (val int) {
	return int(rb.GetUint32())
}

func (rb *ReadBuffer) GetUint64() (val uint64) {
    _ = (*rb)[7]
    val = uint64((*rb)[0])<<56 | uint64((*rb)[1])<<48 | uint64((*rb)[2])<<40 | uint64((*rb)[3])<<32 |
        uint64((*rb)[4])<<24 | uint64((*rb)[5])<<16 | uint64((*rb)[6])<<8 | uint64((*rb)[7])
    *rb = (*rb)[8:]
    return
}

func (rb *ReadBuffer) GetCString() (s string) {
    start, pos := 0, 0
    for (*rb)[pos] != 0 {
        pos++
    }

    l := pos - start
    s = string((*rb)[:l])
    *rb = (*rb)[l+1:]
    return
}

func (rb *ReadBuffer) Get(n int) (bs []byte) {
    _ = (*rb)[n-1]
    bs = (*rb)[:n]
    *rb = (*rb)[n:]
    return
}

func (rb *ReadBuffer) GetCopy(n int) []byte {
    _ = (*rb)[n-1]
    copied := make([]byte, n)
    copy(copied, (*rb)[:n])
    *rb = (*rb)[n:]
    return copied
}

func (rb *ReadBuffer) Remaining() int {
    return len(*rb)
}