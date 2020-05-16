package pixelline

import (
	"fmt"
	"image/color"
)

// getPartCount tells you how many parts there are in a color model.
func getPartCount(m color.Model) int {
	switch m {
	case color.NRGBAModel, color.NRGBA64Model, color.RGBAModel, color.RGBA64Model, color.CMYKModel, color.NYCbCrAModel:
		return 4
	case color.YCbCrModel:
		return 3
	case color.GrayModel, color.Gray16Model, color.AlphaModel, color.Alpha16Model:
		return 1
	default:
		panic("Unknown color model")
	}
}

// getParts decomposes a color into it's parts, converting to float64.
func getParts(c color.Color) []float64 {
	switch v := c.(type) {
	case color.RGBA:
		return []float64{float64(v.R), float64(v.G), float64(v.B), float64(v.A)}
	case color.NRGBA:
		return []float64{float64(v.R), float64(v.G), float64(v.B), float64(v.A)}
	case color.RGBA64:
		return []float64{float64(v.R), float64(v.G), float64(v.B), float64(v.A)}
	case color.NRGBA64:
		return []float64{float64(v.R), float64(v.G), float64(v.B), float64(v.A)}
	case color.CMYK:
		return []float64{float64(v.C), float64(v.M), float64(v.Y), float64(v.K)}
	case color.NYCbCrA:
		return []float64{float64(v.Y), float64(v.Cb), float64(v.Cr), float64(v.A)}
	case color.YCbCr:
		return []float64{float64(v.Y), float64(v.Cb), float64(v.Cr)}
	case color.Gray:
		return []float64{float64(v.Y)}
	case color.Gray16:
		return []float64{float64(v.Y)}
	case color.Alpha:
		return []float64{float64(v.A)}
	case color.Alpha16:
		return []float64{float64(v.A)}
	default:
		panic(fmt.Sprintf("Unknown color type %T", c))
	}
}

// putParts assembles parts for a model into a color.
// Order muse be the same as returned by getParts.
func putParts(m color.Model, parts []float64) color.Color {
	switch m {
	case color.NRGBAModel:
		return color.NRGBA{R: uint8(parts[0]), G: uint8(parts[1]), B: uint8(parts[2]), A: uint8(parts[3])}
	case color.NRGBA64Model:
		return color.NRGBA64{R: uint16(parts[0]), G: uint16(parts[1]), B: uint16(parts[2]), A: uint16(parts[3])}
	case color.RGBAModel:
		return color.RGBA{R: uint8(parts[0]), G: uint8(parts[1]), B: uint8(parts[2]), A: uint8(parts[3])}
	case color.RGBA64Model:
		return color.RGBA64{R: uint16(parts[0]), G: uint16(parts[1]), B: uint16(parts[2]), A: uint16(parts[3])}
	case color.CMYKModel:
		return color.CMYK{C: uint8(parts[0]), M: uint8(parts[1]), Y: uint8(parts[2]), K: uint8(parts[3])}
	case color.NYCbCrAModel:
		return color.NYCbCrA{YCbCr: color.YCbCr{Y: uint8(parts[0]), Cb: uint8(parts[1]), Cr: uint8(parts[2])}, A: uint8(parts[3])}
	case color.YCbCrModel:
		return color.YCbCr{Y: uint8(parts[0]), Cb: uint8(parts[1]), Cr: uint8(parts[2])}
	case color.GrayModel:
		return color.Gray{Y: uint8(parts[0])}
	case color.Gray16Model:
		return color.Gray16{Y: uint16(parts[0])}
	case color.AlphaModel:
		return color.Alpha{A: uint8(parts[0])}
	case color.Alpha16Model:
		return color.Alpha16{A: uint16(parts[0])}
	default:
		panic("Unknown color model")
	}
}
