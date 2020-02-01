package jagex

import (
    "fmt"
    "io"
)

const (
    cacheDescriptorLength   = 6
    cacheBlockHeaderLength  = 8
    cacheBlockContentLength = 512
    cacheBlockLength        = cacheBlockHeaderLength + cacheBlockContentLength
    cacheBlockEOF           = 0
)

type CacheReader struct {
    index  io.ReadSeeker
    blocks io.ReadSeeker
    fileType uint8
    buf      [cacheBlockLength]byte
}

func NewCacheReader(index io.ReadSeeker, blocks io.ReadSeeker, fileType uint8) *CacheReader {
    return &CacheReader{
        index:    index,
        blocks:   blocks,
        fileType: fileType,
    }
}

func (r *CacheReader) Get(fileID uint16) ([]byte, error) {
    fd, err := r.readFileDescriptor(fileID)
    if err != nil {
        return nil, err
    }

    var (
        file  = make([]byte, fd.length)
        nonce = uint16(0)
        bid   = fd.start
    )

    for offset := uint32(0); offset < fd.length; nonce++ {
        if bid == cacheBlockEOF {
            return nil, fmt.Errorf(
                "reached unexpected EOF; write-offset: %d, file-length: %d", offset, fd.length,
            )
        }

        block, err := r.readBlock(bid)
        if err != nil {
            return nil, err
        }

        if block.fid != fileID {
            return nil, fmt.Errorf(
                "file identifier mismatch; expected: %d, found: %d", fileID, block.fid,
            )
        }

        if block.nonce != nonce {
            return nil, fmt.Errorf(
                "file nonce mismatch; expected: %d, found: %d", nonce, block.nonce,
            )
        }

        if block.ft != r.fileType {
            return nil, fmt.Errorf(
                "file type mismatch; expected: %d, found: %d",
                r.fileType,
                block.ft,
            )
        }

        n := fd.length - offset
        if n > cacheBlockContentLength {
            n = cacheBlockContentLength
        }

        copy(
            file[offset:offset+n],
            block.data[:n],
        )

        bid = block.bid
        offset += n
    }

    return file, nil
}

type fileDescriptor struct {
    length uint32
    start  uint32
}

func (r *CacheReader) readFileDescriptor(fileID uint16) (desc fileDescriptor, err error) {
    if _, err = r.index.Seek(int64(fileID*cacheDescriptorLength), io.SeekStart); err != nil {
        return
    }

    if _, err = io.ReadFull(r.index, r.buf[:cacheDescriptorLength]); err != nil {
        return
    }

    rb := ReadBuffer(r.buf[:cacheDescriptorLength])

    return fileDescriptor{
        length: rb.GetUint24(),
        start:  rb.GetUint24(),
    }, nil
}

type cacheBlock struct {
    fid   uint16
    nonce uint16
    bid   uint32
    ft    uint8
    data  []byte
}

func (r *CacheReader) readBlock(blockID uint32) (block cacheBlock, err error) {
    if _, err = r.blocks.Seek(int64(blockID*cacheBlockLength), io.SeekStart); err != nil {
        return
    }

    if _, err = io.ReadFull(r.blocks, r.buf[:cacheBlockLength]); err != nil {
        return
    }

    rb := ReadBuffer(r.buf[:cacheBlockLength])

    return cacheBlock{
        fid:   rb.GetUint16(),
        nonce: rb.GetUint16(),
        bid:   rb.GetUint24(),
        ft:    rb.GetUint8(),
        data:  r.buf[cacheBlockHeaderLength:],
    }, nil
}
