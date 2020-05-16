package pixelline

import (
	"fmt"
	"image/color"
	"testing"
)

func TestGetPartCount(t *testing.T) {
	t.Parallel()

	answerMap := make(map[int][]color.Model, 3)
	answerMap[4] = []color.Model{color.NRGBAModel, color.NRGBA64Model, color.RGBAModel, color.RGBA64Model, color.CMYKModel, color.NYCbCrAModel}
	answerMap[3] = []color.Model{color.YCbCrModel}
	answerMap[1] = []color.Model{color.GrayModel, color.Gray16Model, color.AlphaModel, color.Alpha16Model}

	for count, models := range answerMap {
		for _, model := range models {
			t.Run("getPartCount()", func(t *testing.T) {
				got := getPartCount(model)
				if got != count {
					t.Errorf("Color Model %s returned wrong part count: (got %d expected %d)", model, count, got)
				}
			})
		}
	}
}

func TestGetPartsLength(t *testing.T) {
	t.Parallel()

	answerMap := make(map[int][]color.Color, 3)
	answerMap[4] = []color.Color{color.NRGBA{}, color.NRGBA64{}, color.RGBA{}, color.RGBA64{}, color.CMYK{}, color.NYCbCrA{}}
	answerMap[3] = []color.Color{color.YCbCr{}}
	answerMap[1] = []color.Color{color.Gray{}, color.Gray16{}, color.Alpha{}, color.Alpha16{}}

	for count, colors := range answerMap {
		for _, c := range colors {
			t.Run(fmt.Sprintf("len(getParts(%T))", c), func(t *testing.T) {
				got := len(getParts(c))
				if got != count {
					t.Errorf("Color %s returned the wrong number of parts: (got %d expected %d)", c, got, count)
				}
			})
		}
	}
}

func TestPutPartsCorrectType(t *testing.T) {
	t.Parallel()

	type colorModel struct {
		c color.Color
		m color.Model
	}

	answerMap := make(map[int][]colorModel, 3)
	answerMap[4] = []colorModel{
		{color.NRGBA{}, color.NRGBAModel},
		{color.NRGBA64{}, color.NRGBA64Model},
		{color.RGBA{}, color.RGBAModel},
		{color.RGBA64{}, color.RGBA64Model},
		{color.CMYK{}, color.CMYKModel},
		{color.NYCbCrA{}, color.NYCbCrAModel},
	}
	answerMap[3] = []colorModel{
		{color.YCbCr{}, color.YCbCrModel},
	}
	answerMap[1] = []colorModel{
		{color.Gray{}, color.GrayModel},
		{color.Gray16{}, color.Gray16Model},
		{color.Alpha{}, color.AlphaModel},
		{color.Alpha16{}, color.Alpha16Model},
	}

	for count, colorModels := range answerMap {
		for _, cm := range colorModels {
			t.Run(fmt.Sprintf("putParts(%T).(type)", cm.c), func(t *testing.T) {
				got := putParts(cm.m, make([]float64, count))
				// Compare types without reflection! :D
				// Only because they're both all zeros...
				if got != cm.c {
					t.Errorf("Color model returned the wrong color type: (got %T expected %T)", got, cm.c)
				}
			})
		}
	}
}
