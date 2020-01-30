package jagex

const (
    flagNamed = 0x1
)

type FileIndex struct {
    Revision uint32
}

func DecodeFileIndex(bs []byte) (*FileIndex, error) {
    rb := ReadBuffer(bs)

    version := rb.GetUint8()

    if version >= 6 {
        rb.GetUint32()
    }

    var (
        flags      = rb.GetUint8()
        groupCount = rb.GetUint16()
        groupIDs   = make([]uint16, groupCount)
        id         = uint16(0)
    )

    for i := 0; i < len(groupIDs); i++ {
        id += rb.GetUint16()
        groupIDs[i] = id
    }

    if flags&flagNamed != 0 {

    }

    // Checksum

    // Revision

    //

    return nil, nil
}
