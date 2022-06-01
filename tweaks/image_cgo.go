//go:build cgo
package tweaks

import (
	"io"
	"io/ioutil"
	"math"
	"os"

	"github.com/h2non/bimg"
)

func ResizeImage(infile *io.Reader, outfile *os.File, scale, quality int) {
	buffer, _ := ioutil.ReadAll(*infile)

	img := bimg.NewImage(buffer)
	imgsize, _ := img.Size()

	factor := float64(100 / scale)
	newx1 := int(math.Round(float64(imgsize.Width) / factor))
	newy1 := int(math.Round(float64(imgsize.Height) / factor))

	options := bimg.Options{
		Width: newx1,
		Height: newy1,
		Quality: quality,
	}

	newImage, _ := img.Process(options)

	outfile.Write(newImage)
}

func ImgQuality(infile *io.Reader, outfile *os.File, quality int) {
	options := bimg.Options{
		Quality: quality,
	}

	buffer, _ := ioutil.ReadAll(*infile)

	newImage, _ := bimg.NewImage(buffer).Process(options)

	outfile.Write(newImage)
}