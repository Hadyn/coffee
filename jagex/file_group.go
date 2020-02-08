package jagex

type MultipartFile struct {
    Parts [][]byte
}

func (f *MultipartFile) Collapse() []byte {
    collapsed := make([]byte, 0, f.Length())
    for _, bs := range f.Parts {
        collapsed = append(collapsed, bs...)
    }
    return collapsed
}

func (f *MultipartFile) Length() (length int) {
    for _, bs := range f.Parts {
        length += len(bs)
    }
    return
}

func DecodeFileGroup(bs []byte, size int) []*MultipartFile {
    if size < 2 {
        copied := make([]byte, len(bs))
        copy(copied, bs)

        return []*MultipartFile{{Parts: [][]byte{copied}}}
    }

    rb := ReadBuffer(bs[len(bs)-1:])

    var (
        parts         = rb.GetUint8AsInt()
        lengthsOffset = len(bs) - 4*parts - 1
        lengths       = make([][]int, size)
    )

    rb = bs[lengthsOffset : len(bs)-1]
    for part := 0; part < parts; part++ {
        for fileID := 0; fileID < size; fileID++ {
            lengths[fileID][part] = rb.GetUint32AsInt()
        }
    }

    rb = bs[0:lengthsOffset]

    files := make([]*MultipartFile, size)
    for fileID := 0; fileID < size; fileID++ {
        file := &MultipartFile{Parts: make([][]byte, parts)}
        for part := 0; part < parts; part++ {
            file.Parts[part] = rb.GetCopy(lengths[fileID][part])
        }
        files[fileID] = file
    }

    return files
}
