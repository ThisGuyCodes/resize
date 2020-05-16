package resize

import (
	"image"
	"image/color"
	"sync"

	"github.com/thisguycodes/resize/pixelline"
)

// concurrency limits the number of active goroutines for any one Resize
// operation.
const concurrency = 4

// Resize resizes an image.
// New image always uses the YCbCr color.Model.
func Resize(src image.Image, newSize image.Rectangle) image.Image {
	imgWidth := src.Bounds().Dx()
	imgHeight := src.Bounds().Dy()

	newImg := image.NewYCbCr(newSize, image.YCbCrSubsampleRatio444)
	newImgWidth := newImg.Bounds().Dx()
	newImgHeight := newImg.Bounds().Dy()

	verticleLines := make([]pixelline.PixelLine, imgWidth)

	wg := &sync.WaitGroup{}

	wg.Add(src.Bounds().Dx())
	limit := make(chan struct{}, concurrency)

	for imgX := src.Bounds().Min.X; imgX < src.Bounds().Max.X; imgX++ {
		limit <- struct{}{}
		go func(imgX int) {
			defer wg.Done()
			defer func() { <-limit }()

			normalizedX := imgX - src.Bounds().Min.X + newImg.Bounds().Min.X

			verticleLine := pixelline.New(imgHeight, src.ColorModel())

			for imgY := src.Bounds().Min.Y; imgY < src.Bounds().Max.Y; imgY++ {
				normalizedY := imgY - src.Bounds().Min.Y + src.Bounds().Min.Y
				verticleLine.Set(normalizedY, src.At(imgX, imgY))
			}

			verticleLines[normalizedX] = verticleLine.Resize(newImgHeight)
		}(imgX)
	}
	wg.Wait()

	wg.Add(newImg.Bounds().Dy())
	for y := newImg.Bounds().Min.Y; y < newImg.Bounds().Max.Y; y++ {
		limit <- struct{}{}
		go func(y int) {
			defer wg.Done()
			defer func() { <-limit }()

			horizontalLine := pixelline.New(imgWidth, src.ColorModel())
			for x, verticleLine := range verticleLines {
				horizontalLine.SetParts(x, verticleLine.AtParts(y))
			}

			horizontalLine = horizontalLine.Resize(newImgWidth)

			for x := newImg.Bounds().Min.X; x < newImg.Bounds().Max.X; x++ {
				thisPixel := color.YCbCrModel.Convert(horizontalLine.At(x)).(color.YCbCr)
				yi := newImg.YOffset(x, y)
				ci := newImg.COffset(x, y)
				newImg.Y[yi] = thisPixel.Y
				newImg.Cb[ci] = thisPixel.Cb
				newImg.Cr[ci] = thisPixel.Cr
			}
		}(y)
	}
	wg.Wait()

	return newImg
}
