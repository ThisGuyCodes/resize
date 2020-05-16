package main

import (
	"image"
	"image/jpeg"
	"log"
	"math"
	"os"

	// Support more image types!
	_ "image/jpeg"
	_ "image/png"

	"github.com/thisguycodes/resize"
)

func halfSize(img image.Image) *image.YCbCr {
	newRect := image.Rectangle{
		Max: image.Point{
			X: img.Bounds().Dx() / 2,
			Y: img.Bounds().Dy() / 2,
		},
	}
	return image.NewYCbCr(newRect, image.YCbCrSubsampleRatio444)
}

func main() {
	img, _, err := image.Decode(os.Stdin)

	if err != nil {
		log.Fatal(err)
	}

	ratio := float64(img.Bounds().Dx()) / float64(img.Bounds().Dy())

	tooWide := img.Bounds().Dx() > img.Bounds().Dy()
	var newWidth, newHeight int

	if tooWide {
		newWidth = 1024
		newHeight = int(math.Round(float64(newWidth) / ratio))
	} else {
		newHeight = 1024
		newWidth = int(math.Round(float64(newHeight) * ratio))
	}

	newImgSize := image.Rectangle{
		Max: image.Pt(newWidth, newHeight),
	}

	newImg := resize.Resize(img, newImgSize)

	err = jpeg.Encode(os.Stdout, newImg, &jpeg.Options{Quality: 90})
	if err != nil {
		log.Fatal(err)
	}
}
