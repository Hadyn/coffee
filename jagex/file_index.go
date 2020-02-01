package jagex

import (
    "fmt"
    "github.com/hadyn/coffee/jagex/dbj2"
)

type FileIndexFlags uint8

const (
    IndexFlagNamed FileIndexFlags = 0x1
)

const (
    oldestIndexFormat uint8 = 5
    latestIndexFormat uint8 = 6
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
    if format < oldestIndexFormat || format > latestIndexFormat {
        return nil, fmt.Errorf("format is not supported: %d", format)
    }

    if format >= 6 {
        fi.Revision = rb.GetUint32()
    }

    var (
        flags      = FileIndexFlags(rb.GetUint8())
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

    fi.Groups = make([]*FileGroupEntry, maximumGroupID+1)
    for _, groupID := range groupIDs {
        fi.Groups[groupID] = &FileGroupEntry{}
    }

    if flags&IndexFlagNamed != 0 {
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

    for i := range groupIDs {
        fileCount := rb.GetUint16()
        fileIDs[i] = make([]uint16, fileCount)
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

    for i, groupID := range groupIDs {
        group := fi.Groups[groupID]
        group.Files = make([]*FileEntry, maximumFileIDs[i]+1)

        for _, fileID := range fileIDs[i] {
            group.Files[fileID] = &FileEntry{}
        }
    }

    if flags&IndexFlagNamed != 0 {
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

func (fi *FileIndex) LookupGroupID(groupName string) (groupID int, found bool) {
   	var (
   		group *FileGroupEntry
   		hash = dbj2.Sum([]byte(groupName))
	)

    for groupID, group = range fi.Groups {
        if group.NameHash == hash {
        	found = true
            return
        }
    }

    groupID, found = 0, false
    return
}

func (fi *FileIndex) LookupGroup(groupName string) *FileGroupEntry {
	groupID, found := fi.LookupGroupID(groupName)
	if !found {
		return nil
	}
	return fi.Groups[groupID]
}

func (fi *FileIndex) LookupFileID(groupName, fileName string) (groupID int, fileID int, found bool) {
    groupID, found = fi.LookupGroupID(groupName)
    if !found {
        return
    }

    fileID, found = fi.Groups[groupID].LookupFileID(fileName)
	return
}

func (e *FileGroupEntry) LookupFileID(fileName string) (fileID int, found bool) {
	var (
		file *FileEntry
		hash = dbj2.Sum([]byte(fileName))
	)

	for fileID, file = range e.Files {
		if file.NameHash == hash {
			found = true
			return
		}
	}

	fileID, found = 0, false
	return
}
