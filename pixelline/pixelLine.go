package pixelline

import (
	"fmt"
	"image/color"
	"math"

	"github.com/cnkei/gospline"
)

// New initializes a PixelLine of the given length and color model.
// It is an error to create a line of <= 0 length.
func New(l int, model color.Model) PixelLine {
	if l <= 1 {
		panic("Lines of 1 or less length are invalid")
	}
	partCount := getPartCount(model)
	parts := make([][]float64, partCount)
	for i := range parts {
		parts[i] = make([]float64, l)
	}
	return PixelLine{
		len:   l,
		model: model,
		parts: parts,
	}
}

// PixelLine is a line of pixels. Most noteably it is resizeable.
type PixelLine struct {
	len   int
	model color.Model
	parts [][]float64
}

// Set makes the color at pos equal to c.
func (p *PixelLine) Set(pos int, c color.Color) {
	c = p.model.Convert(c)
	parts := getParts(c)
	p.SetParts(pos, parts)
}

// SetParts lets you backdoor a bit and set the underlying for a given position.
func (p *PixelLine) SetParts(pos int, parts []float64) {
	for i, part := range parts {
		p.parts[i][pos] = part
	}
}

// At converts and returns the color at the given position in the line.
// Safe for concurrent use of different positions.
func (p *PixelLine) At(pos int) color.Color {
	return putParts(p.model, p.AtParts(pos))
}

// AtParts lets you backdoor a bit and retrive the underlying parts for a given position.
func (p *PixelLine) AtParts(pos int) []float64 {
	parts := make([]float64, len(p.parts))
	for i := range parts {
		parts[i] = p.parts[i][pos]
	}
	return parts
}

// Len gives the current length of the line.
func (p *PixelLine) Len() int {
	return p.len
}

// Resize returns a new PixelLine that is a resized version of this one.
func (p *PixelLine) Resize(newLen int) PixelLine {
	if newLen <= 1 {
		panic("Lines of 1 or less length are invalid")
	}
	newSpline := gospline.NewMonotoneSpline
	locs := genLocs(p.len)

	end := float64(locs[len(locs)-1])
	step := end / float64(newLen-1)

	// Because float rounding can screw us.
	// This is based on the observed math inside cubic.Range in gospline
	//
	// Found with fuzzing!
	for int(end/step) < newLen-1 {
		step = math.Nextafter(step, 0)
	}

	newParts := make([][]float64, len(p.parts))
	fmt.Printf("from %d to %d\n", p.len, newLen)
	for i, part := range p.parts {
		newParts[i] = newSpline(locs, part).Range(0, end, step)
	}
	return PixelLine{
		len:   newLen,
		model: p.model,
		parts: newParts,
	}
}

// genLocs creates a slice of floats pre-filled with their indicies.
// Panics on l < 0.
func genLocs(l int) []float64 {
	ret := make([]float64, l)
	for i := range ret {
		ret[i] = float64(i)
	}
	return ret
}
