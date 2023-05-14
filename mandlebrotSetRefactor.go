package main

import (
	"os"
	"sync"

	"math/cmplx"

	"image"
	"image/color"
	"image/png"
)

func main() {
	width, height := 1_000, 1_000
	maxIterations := 10_000

	finalImg := mandlebrot(width, height, maxIterations)
	f, _ := os.Create("mandlebrot.png")
	png.Encode(f, finalImg)
}

func mandlebrot(width int, height int, maxIterations int) (finalImg *image.RGBA) {
	xNormal, yNormal := 2.5/float64(width), 2.5/float64(height)
	var upLeft, lowRight = image.Point{0, 0}, image.Point{width, height}
	var img = image.NewRGBA(image.Rectangle{upLeft, lowRight})
	var wg sync.WaitGroup

	draw := func(x int, y int, xNormal float64, yNormal float64, maxIterations int, calcEndIter func(float64, float64, int) int, img **image.RGBA, wg *sync.WaitGroup) {
		endIterationScaled := uint8(calcEndIter(float64(x)*xNormal-2, float64(y)*yNormal-1.12, maxIterations) * 255 / maxIterations)

		(*img).Set(x, y, color.RGBA{0, 100, endIterationScaled, 0xff})
		(*wg).Done()
	}
	calcEndIter := func(xScaled float64, yScaled float64, maxIterations int) (endIteration int) {
		var c = complex(xScaled, yScaled)
		var z = complex(0, 0)

		var iteration int = 0
		for cmplx.Abs(z) < 2 && iteration < maxIterations {
			z = cmplx.Pow(z, 2) + c
			iteration++
		}
		return iteration
	}
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			wg.Add(1)
			go draw(x, y, xNormal, yNormal, maxIterations, calcEndIter, &img, &wg)
		}
	}
	wg.Wait()

	return img
}
