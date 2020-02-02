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
    Revision uint32            `json:"revision"`
    Groups   []*FileGroupEntry `json:"groups"`
}

type FileGroupEntry struct {
    NameHash uint32       `json:"name"`
    Checksum uint32       `json:"checksum"`
    Revision uint32       `json:"revision"`
    Files    []*FileEntry `json:"files"`
}

type FileEntry struct {
    NameHash uint32 `json:"name"`
}

type NamedEntryIndex struct {
    groupLookup map[uint32]int
    fileLookup  map[int]map[uint32]int
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
        group.Revision = rb.GetUint32()
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

func (fi *FileIndex) Size() (size int) {
    for _, group := range fi.Groups {
        if group == nil {
            continue
        }
        size++
    }
    return
}

func (fi *FileIndex) FindGroupID(groupName string) (groupID int, found bool) {
    var (
        group *FileGroupEntry
        hash  = dbj2.Sum([]byte(groupName))
    )

    for groupID, group = range fi.Groups {
        if group == nil {
            continue
        }

        if group.NameHash == hash {
            found = true
            return
        }
    }

    groupID, found = 0, false
    return
}

func (fi *FileIndex) FindGroup(groupName string) *FileGroupEntry {
    groupID, found := fi.FindGroupID(groupName)
    if !found {
        return nil
    }
    return fi.Groups[groupID]
}

func (fi *FileIndex) FindFileID(groupName, fileName string) (groupID int, fileID int, found bool) {
    groupID, found = fi.FindGroupID(groupName)
    if !found {
        return
    }

    fileID, found = fi.Groups[groupID].FindFileID(fileName)
    return
}

func (fi *FileIndex) NamedIndex() *NamedEntryIndex {
    var (
        size  = fi.Size()
        index = &NamedEntryIndex{
            groupLookup: make(map[uint32]int, size),
            fileLookup:  make(map[int]map[uint32]int, size),
        }
    )

    for groupID, group := range fi.Groups {
        if group == nil {
            continue
        }

        index.groupLookup[group.NameHash] = groupID

        index.fileLookup[groupID] = make(map[uint32]int, group.Size())
        for fileID, file := range group.Files {
            if file == nil {
                continue
            }

            index.fileLookup[groupID][file.NameHash] = fileID
        }
    }

    return index
}

func (e *FileGroupEntry) Size() (size int) {
    for _, file := range e.Files {
        if file == nil {
            continue
        }
        size++
    }
    return
}

func (e *FileGroupEntry) FindFileID(fileName string) (fileID int, found bool) {
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

func (e *NamedEntryIndex) LookupGroupID(groupName string) (groupID int, exists bool) {
    groupID, exists = e.groupLookup[dbj2.Sum([]byte(groupName))]
    return
}

func (e *NamedEntryIndex) LookupFileID(groupName, fileName string) (groupID int, fileID int, exists bool) {
    groupID, exists = e.LookupGroupID(groupName)
    if !exists {
        return
    }

    fileID, exists = e.fileLookup[groupID][dbj2.Sum([]byte(fileName))]
    if exists {
        return
    }

    groupID, fileID, exists = 0, 0, false
    return
}
