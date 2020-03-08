package sprite

import (
    "fmt"
    "github.com/hadyn/coffee/jagex"
    "image"
    "image/color"
)

const (
    TransparentPixel uint32 = 0x000000
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

    for x := 0; x < s.PackedWidth; x++ {
        for y := 0; y < s.PackedHeight; y++ {
            var (
                src   = s.Colors[s.Index[x+y*s.PackedWidth]]
                pixel color.RGBA
            )

            if src != TransparentPixel {
                pixel.R, pixel.G, pixel.B = uint8(src >> 16 & 0xff), uint8(src >> 8 & 0xff), uint8(src & 0xff)
                pixel.A = 0xff
            }

            img.Set(x+s.OffsetX, y+s.OffsetY, pixel)
        }
    }

    return img
}

type Group struct {
    Width        int
    Height       int
    OffsetX      []int
    OffsetY      []int
    PackedWidth  []int
    PackedHeight []int
    Colors       []uint32
    Index        [][]uint8
}

func DecodeGroup(bs []byte) *Group {
    rb := jagex.ReadBuffer(bs[len(bs)-2:])

    count := rb.GetUint16AsInt()

    rb = bs[len(bs)-7-count*8:]

    group := &Group{
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
        group.OffsetX[i] = rb.GetUint16AsInt()
    }

    for i := 0; i < count; i++ {
        group.OffsetY[i] = rb.GetUint16AsInt()
    }

    for i := 0; i < count; i++ {
        group.PackedWidth[i] = rb.GetUint16AsInt()
    }

    for i := 0; i < count; i++ {
        group.PackedHeight[i] = rb.GetUint16AsInt()
    }

    rb = bs[len(bs)-7-count*8-(len(group.Colors)-1)*3:]

    for i := 1; i < len(group.Colors); i++ {
        group.Colors[i] = rb.GetUint24()

        if group.Colors[i] == TransparentPixel {
            group.Colors[i] = 1
        }
    }

    rb = bs[0:]

    for i := 0; i < count; i++ {
        pw, ph := group.PackedWidth[i], group.PackedHeight[i]
        size := pw * ph

        group.Index[i] = make([]uint8, size)

        encoding := PixelEncoding(rb.GetUint8())

        switch encoding {
        case HorizontalEncoding:
            for j := 0; j < size; j++ {
                group.Index[i][j] = rb.GetUint8()
            }
        case VerticalEncoding:
            for x := 0; x < pw; x++ {
                for y := 0; y < ph; y++ {
                    group.Index[i][x+pw*y] = rb.GetUint8()
                }
            }
        default:
            panic(fmt.Sprintf("Unrecognized sprite encoding: %d", encoding))
        }
    }

    return group
}

func (s *Group) Get(n int) *Sprite {
    return &Sprite{
        Width:        s.Width,
        Height:       s.Height,
        OffsetX:      s.OffsetX[n],
        OffsetY:      s.OffsetY[n],
        PackedWidth:  s.PackedWidth[n],
        PackedHeight: s.PackedHeight[n],
        Colors:       s.Colors,
        Index:        s.Index[n],
    }
}
