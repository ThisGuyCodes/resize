package pixelline

import (
	"fmt"
	"image/color"
	"testing"

	fuzz "github.com/google/gofuzz"
)

func TestGenLocs(t *testing.T) {
	t.Parallel()

	fuzzer := fuzz.New()

	for i := 0; i < 1000; i++ {
		var num int
		fuzzer.Fuzz(&num)
		// limit the size because memory
		num = num % 1000
		t.Run(fmt.Sprintf("genLocs(%d)", num), func(t *testing.T) {
			if num < 0 {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("genLocs(%d) did not panic", num)
					}
				}()
			}
			got := genLocs(num)
			if len(got) != num {
				t.Errorf("genLocs() returned the wrong length: (got %d expected %d)", len(got), num)
			}
		})
	}
}

func TestPixelLineLen(t *testing.T) {
	t.Parallel()

	fuzzer := fuzz.New()

	models := []color.Model{
		color.NRGBAModel,
		color.NRGBA64Model,
		color.RGBAModel,
		color.RGBA64Model,
		color.CMYKModel,
		color.NYCbCrAModel,
		color.YCbCrModel,
		color.GrayModel,
		color.Gray16Model,
		color.AlphaModel,
		color.Alpha16Model,
	}

	for _, model := range models {
		for i := 0; i < 100; i++ {
			var l int
			fuzzer.Fuzz(&l)
			l = l % 100000

			t.Run(fmt.Sprintf("New(%d, %v).Len()", l, model), func(t *testing.T) {
				if l < 0 {
					defer func() {
						if r := recover(); r == nil {
							t.Errorf("New(%d, %v) did not panic", l, model)
						}
					}()
				}
				line := New(l, model)
				got := line.Len()
				if line.Len() != l {
					t.Errorf("Got unexpected length: (got %d expected %d)", l, got)
				}
			})
		}
	}
}

func TestPixelLineSetAtParts(t *testing.T) {
	t.Parallel()

	fuzzer := fuzz.New()

	answerMap := make(map[int][]color.Model, 3)
	answerMap[4] = []color.Model{color.NRGBAModel, color.NRGBA64Model, color.RGBAModel, color.RGBA64Model, color.CMYKModel, color.NYCbCrAModel}
	answerMap[3] = []color.Model{color.YCbCrModel}
	answerMap[1] = []color.Model{color.GrayModel, color.Gray16Model, color.AlphaModel, color.Alpha16Model}

	for _, models := range answerMap {
		for _, m := range models {
			line := New(10, m)
			for i := 0; i < 100; i++ {
				var pos int
				fuzzer.Fuzz(&pos)
				// limit possible values so we actually test things...
				pos = pos % (10 + 1)

				parts := make([]float64, getPartCount(m))
				for i := range parts {
					fuzzer.Fuzz(&parts[i])
				}

				t.Run(fmt.Sprintf("PixelLine.SetParts(%d, %#v)", pos, parts), func(t *testing.T) {
					if pos < 0 || pos >= 10 {
						defer func() {
							if r := recover(); r == nil {
								t.Errorf("PixelLine.SetParts(%d, %#v) did not panic", pos, parts)
							}
						}()
					}
					line.SetParts(pos, parts)
					gotParts := line.AtParts(pos)
					for pLoc, part := range gotParts {
						if parts[pLoc] != part {
							t.Errorf("PixelLine.AtParts(%d) returned unexpected value: (got %#v expected %#v)", pos, gotParts, parts)
						}
					}
				})
			}
		}
	}
}

func TestPixelLineSetAt(t *testing.T) {
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

	for _, colorModels := range answerMap {
		for _, cm := range colorModels {
			t.Run(fmt.Sprintf("PixelLine.Set(0, %T)", cm.c), func(t *testing.T) {
				line := New(2, cm.m)
				line.Set(0, cm.c)
				got := line.At(0)

				if cm.c != got {
					t.Fail()
				}
			})
		}

	}
}

func TestPixelLineResize(t *testing.T) {
	t.Parallel()

	fuzzer := fuzz.New()
	m := color.GrayModel
	for i := 0; i < 1000; i++ {
		var l, nl int
		fuzzer.Fuzz(&l)
		l = l % 100
		if l < 0 {
			l = l * -1
		}
		if l <= 1 {
			l = 2
		}
		fuzzer.Fuzz(&nl)
		nl = nl % 100
		if nl < 0 {
			nl = nl * -1
		}
		if nl <= 1 {
			nl = 2
		}
		t.Run(fmt.Sprintf("New(%d, ##).Resize(%d)", l, nl), func(t *testing.T) {
			t.Logf("%d, %d", l, nl)
			line := New(l, m)
			line = line.Resize(nl)

			if line.Len() != nl {
				t.Errorf("line didn't resize to the correct length: (got %d expected %d)", line.Len(), nl)
			}
		})
	}
}

func TestPixelLineResizeAlso(t *testing.T) {

	fuzzer := fuzz.New()
	var l int
	fuzzer.Fuzz(&l)
	l = l % 100000
	if l < 0 {
		l = l * -1
	}

	line := New(1000, color.YCbCrModel)

	for i := 0; i < 1000; i++ {
		var l int
		fuzzer.Fuzz(&l)
		l = l % 100
		if l < 0 {
			l = l * -1
		}
		if l <= 1 {
			l = 2
		}
		line = line.Resize(l)
	}
}
