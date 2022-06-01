//go:build !cgo
package tweaks

import (
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"os"

	x_draw "golang.org/x/image/draw"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

func ResizeImage(infile *io.ReadCloser, outfile *os.File, scale, quality int) {
	src, format, _ := image.Decode(*infile)

	factor := float64(100 / scale)
	newx1 := int(math.Round(float64(src.Bounds().Max.X)/factor))
	newy1 := int(math.Round(float64(src.Bounds().Max.Y)/factor))

	// Set the expected size that you want:
	dst := image.NewNRGBA(image.Rect(0, 0, newx1, newy1))

	// Scale:
	x_draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)

	// Encode to `output`:

	switch format {
	case "png":
		png.Encode(outfile, dst)
	case "jpeg":
		jpeg.Encode(outfile, dst, &jpeg.Options{Quality: quality})
	case "gif":
		gif.Encode(outfile, dst, &gif.Options{})
	}
}

func ImgQuality(infile *io.ReadCloser, outfile *os.File, quality int) {
	src, format, _ := image.Decode(*infile)

	if format == "jpeg" {
		options := jpeg.Options{Quality: quality}

		jpeg.Encode(outfile, src, &options)
	}
}