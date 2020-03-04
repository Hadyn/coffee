package jagex

// ReadBuffer provides reading of various binary types from byte slice. After reading a value,
// the underlying slice will be resliced so that unread bytes are next to be read.
type ReadBuffer []byte

// Remaining is an alias for the length of the underlying slice.
func (rb *ReadBuffer) Remaining() int {
	return len(*rb)
}

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

// GetCString returns a null terminated string.
// TODO(hadyn): Runescape uses a specific character set, need to implement that.
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

// Get returns a slice of the next N bytes and reslices the underlying buffer to be a view of the
// bytes after the read bytes. This method should be carefully used because changes to the
// elements in the returned slice will appear in this buffer and elsewhere.
func (rb *ReadBuffer) Get(n int) (bs []byte) {
	_ = (*rb)[n-1]
	bs = (*rb)[:n]
	*rb = (*rb)[n:]
	return
}

// Get returns a copy of the next n bytes and reslices the underlying buffer to be a view of
// the bytes there after. This method is similar to Get except it returns a copy of the bytes
// which makes it safe to do modifications to the elements.
func (rb *ReadBuffer) Copy(n int) []byte {
	if n == 0 {
		return []byte{}
	}

	_ = (*rb)[n-1]
	copied := make([]byte, n)
	copy(copied, (*rb)[:n])
	*rb = (*rb)[n:]
	return copied
}
