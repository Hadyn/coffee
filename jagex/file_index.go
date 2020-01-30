package jagex

const (
	flagNamed = 0x1
)

type FileIndex struct {
	Revision uint32
	Groups   []*FileGroupEntry
}

type FileGroupEntry struct {
	NameHash uint32
	Checksum uint32
	Revision uint16
	Files    []*FileEntry
}

type FileEntry struct {
	NameHash uint32
}

func DecodeFileIndex(bs []byte) (*FileIndex, error) {
	var (
		fi = &FileIndex{}
		rb = ReadBuffer(bs)
	)

	format := rb.GetUint8()
	if format >= 6 {
		fi.Revision = rb.GetUint32()
	}

	var (
		flags      = rb.GetUint8()
		groupCount = rb.GetUint16()
	)

	var (
		groupIDs       = make([]uint16, groupCount)
		currentGroupID = uint16(0)
		maximumGroupID = uint16(0)
	)

	for i := range groupIDs {
		currentGroupID += rb.GetUint16()
		if currentGroupID > maximumGroupID {
			maximumGroupID = currentGroupID
		}
		groupIDs[i] = currentGroupID
	}

	if flags&flagNamed != 0 {
		for _, groupID := range groupIDs {
			group := fi.Groups[groupID]
			group.NameHash = rb.GetUint32()
		}
	}

	for _, groupID := range groupIDs {
		group := fi.Groups[groupID]
		group.Checksum = rb.GetUint32()
	}

	for _, groupID := range groupIDs {
		group := fi.Groups[groupID]
		group.Revision = rb.GetUint16()
	}

	var (
		fileIDs        = make([][]uint16, groupCount)
		maximumFileIDs = make([]uint16, groupCount)
	)

	for _, groupID := range groupIDs {
		fileCount := rb.GetUint16()
		fileIDs[groupID] = make([]uint16, fileCount)
	}

	for i := range groupIDs {
		currentFileID := uint16(0)
		for j := range fileIDs[i] {
			currentFileID += rb.GetUint16()

			if currentFileID > maximumFileIDs[i] {
				maximumFileIDs[i] = currentFileID
			}

			fileIDs[i][j] = currentFileID
		}
	}

	if flags&flagNamed != 0 {
		for i, groupID := range groupIDs {
			for _, fileID := range fileIDs[i] {
				var (
					group = fi.Groups[groupID]
					file  = group.Files[fileID]
				)

				file.NameHash = rb.GetUint32()
			}
		}
	}

	return fi, nil
}
