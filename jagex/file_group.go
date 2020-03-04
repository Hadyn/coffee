package jagex

type MultipartFile struct {
	Parts [][]byte
}

func (f *MultipartFile) Collapse() []byte {
	if len(f.Parts) <= 1 {
		return f.Parts[0]
	}

	collapsed := make([]byte, 0, f.CollapsedLength())
	for _, bs := range f.Parts {
		collapsed = append(collapsed, bs...)
	}

	return collapsed
}

func (f *MultipartFile) CollapsedLength() (length int) {
	for _, bs := range f.Parts {
		length += len(bs)
	}
	return
}

func DecodeFileGroup(bs []byte, size int) []*MultipartFile {
	if size == 0 {
		return []*MultipartFile{}
	}

	if size == 1 {
		copied := make([]byte, len(bs))
		copy(copied, bs)

		return []*MultipartFile{{Parts: [][]byte{copied}}}
	}

	rb := ReadBuffer(bs[len(bs)-1:])

	var (
		parts         = rb.GetUint8AsInt()
		lengthsOffset = len(bs) - 4*parts*size - 1
		lengths       = make([][]int, size)
	)

	rb = bs[lengthsOffset : len(bs)-1]
	for part := 0; part < parts; part++ {
		for fileID := 0; fileID < size; fileID++ {
			if lengths[fileID] == nil {
				lengths[fileID] = make([]int, parts)
			}
			lengths[fileID][part] = rb.GetUint32AsInt()
		}
	}

	rb = bs[0:lengthsOffset]

	files := make([]*MultipartFile, size)
	for part := 0; part < parts; part++ {
		for fileID := 0; fileID < size; fileID++ {
			if files[fileID] == nil {
				files[fileID] = &MultipartFile{Parts: make([][]byte, parts)}
			}
			files[fileID].Parts[part] = rb.Copy(lengths[fileID][part])
		}
	}

	return files
}
