package sprite

import (
    "fmt"
    "github.com/hadyn/coffee/jagex"
    "image"
    "image/color"
)

type PixelEncoding uint8

const (
    HorizontalEncoding PixelEncoding = 0
    VerticalEncoding   PixelEncoding = 1
)

type Sprite struct {
    Width        int
    Height       int
    OffsetX      int
    OffsetY      int
    PackedWidth  int
    PackedHeight int
    Colors       []uint32
    Index        []uint8
}

func (s *Sprite) ToImage() image.Image {
    img := image.NewRGBA(image.Rect(
        s.OffsetX,
        s.OffsetY,
        s.OffsetX+s.PackedWidth,
        s.OffsetY+s.PackedHeight,
    ))

    for x := s.OffsetX; x < s.OffsetX+s.PackedWidth; x++ {
        for y := s.OffsetY; y < s.OffsetY+s.PackedHeight; y++ {
            var (
                src   = s.Colors[s.Index[x+y*s.PackedHeight]]
                pixel color.RGBA
            )

            if src != 0x000000 {
                pixel.R = uint8(src >> 16 & 0xff)
                pixel.G = uint8(src >> 8 & 0xff)
                pixel.B = uint8(src & 0xff)
                pixel.A = 0xff
            }

            img.Set(x, y, pixel)
        }
    }

    return img
}

type Sheet struct {
    Width        int
    Height       int
    OffsetX      []int
    OffsetY      []int
    PackedWidth  []int
    PackedHeight []int
    Colors       []uint32
    Index        [][]uint8
}

func DecodeSheet(bs []byte) *Sheet {
    rb := jagex.ReadBuffer(bs[len(bs)-2:])

    count := rb.GetUint16AsInt()

    rb = bs[len(bs)-7-count*8:]

    sheet := &Sheet{
        Width:        rb.GetUint16AsInt(),
        Height:       rb.GetUint16AsInt(),
        OffsetX:      make([]int, count),
        OffsetY:      make([]int, count),
        PackedWidth:  make([]int, count),
        PackedHeight: make([]int, count),
        Colors:       make([]uint32, rb.GetUint8AsInt()+1),
        Index:        make([][]uint8, count),
    }

    for i := 0; i < count; i++ {
        sheet.OffsetX[i] = rb.GetUint16AsInt()
    }

    for i := 0; i < count; i++ {
        sheet.OffsetY[i] = rb.GetUint16AsInt()
    }

    for i := 0; i < count; i++ {
        sheet.PackedWidth[i] = rb.GetUint16AsInt()
    }

    for i := 0; i < count; i++ {
        sheet.PackedHeight[i] = rb.GetUint16AsInt()
    }

    rb = bs[len(bs)-7-count*8-(len(sheet.Colors)-1)*3:]

    for i := 1; i < len(sheet.Colors); i++ {
        sheet.Colors[i] = rb.GetUint24()

        if sheet.Colors[i] == 0x000000 {
            sheet.Colors[i] = 0x000001
        }
    }

    rb = bs[0:]

    for i := 0; i < count; i++ {
        pw, ph := sheet.PackedWidth[i], sheet.PackedHeight[i]
        size := pw * ph

        sheet.Index[i] = make([]uint8, size)

        encoding := PixelEncoding(rb.GetUint8())

        switch encoding {
        case HorizontalEncoding:
            for j := 0; j < size; j++ {
                sheet.Index[i][j] = rb.GetUint8()
            }
        case VerticalEncoding:
            for x := 0; x < pw; x++ {
                for y := 0; y < ph; y++ {
                    sheet.Index[i][x+pw*y] = rb.GetUint8()
                }
            }
        default:
            panic(fmt.Sprintf("Unrecognized sprite encoding: %d", encoding))
        }
    }

    return sheet
}
