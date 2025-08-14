package image

import (
	"fmt"
	"sync"

	"github.com/sendelivery/go-trace-rays/internal/color"
	"github.com/sendelivery/go-trace-rays/internal/interval"
)

type PixelCoord struct {
	x, y int
}

func (pc PixelCoord) X() int {
	return pc.x
}

func (pc PixelCoord) Y() int {
	return pc.y
}

func NewPixelCoord(x, y int) PixelCoord {
	return PixelCoord{x, y}
}

type Chunk struct {
	start, end PixelCoord
}

func (c Chunk) Start() PixelCoord {
	return c.start
}

func (c Chunk) End() PixelCoord {
	return c.end
}

func NewChunk(start, end PixelCoord) Chunk {
	return Chunk{start, end}
}

type Image struct {
	width  interval.Interval
	height interval.Interval
	m      sync.Map
}

func New(width, height int) Image {
	return Image{
		width:  interval.New(0, float64(width-1)),
		height: interval.New(0, float64(height-1)),
	}
}

func (i *Image) Add(pc PixelCoord, col color.Color) error {
	if !i.height.Contains(float64(pc.y)) {
		return fmt.Errorf("y %d is outside of image bounds", pc.y)
	}
	if !i.width.Contains(float64(pc.x)) {
		return fmt.Errorf("x %d is outside of image bounds", pc.x)
	}

	i.m.Store(pc, col)
	return nil
}

func (i *Image) Get(pc PixelCoord) (color.Color, bool) {
	px, ok := i.m.Load(pc)
	if !ok {
		return color.Black, false
	}

	pxc, ok := px.(color.Color)
	return pxc, ok
}
