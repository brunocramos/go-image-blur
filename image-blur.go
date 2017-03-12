package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"
)

type Pixel struct {
	R int
	G int
	B int
	A int
}

// Blur mask size
const MASK_SIZE = 3

// Fill a Pixel struct with provided rgba
func getImagePixelFromRGBA(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}

// Get matrix of pixels from image
func getImagePixels(file io.Reader) ([][]Pixel, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var pixelMatrix [][]Pixel
	for i := 0; i < height; i++ {
		var row []Pixel
		for j := 0; j < width; j++ {
			row = append(row, getImagePixelFromRGBA(img.At(i, j).RGBA()))
		}
		pixelMatrix = append(pixelMatrix, row)
	}

	return pixelMatrix, nil
}

// Process blur mask for given pixel
func processBlurMask(i, j int, pixels [][]Pixel) Pixel {
	offset := 0

	if MASK_SIZE%2 == 0 {
		offset = MASK_SIZE / 2
	} else {
		offset = (MASK_SIZE - 1) / 2
	}

	x := i - offset
	y := j - offset
	var averagePixel Pixel
	total := 0

	for a := 0; a < MASK_SIZE; a++ {
		for b := 0; b < MASK_SIZE; b++ {
			if x+a >= 0 && x+a < len(pixels) {
				if y+b >= 0 && y+b < len(pixels[x+a]) {
					averagePixel.R += pixels[x+a][y+b].R
					averagePixel.G += pixels[x+a][y+b].G
					averagePixel.B += pixels[x+a][y+b].B
					total++
				}
			}
		}
	}

	averagePixel.R = averagePixel.R / total
	averagePixel.G = averagePixel.G / total
	averagePixel.B = averagePixel.B / total
	averagePixel.A = pixels[i][j].A
	return averagePixel
}

// Apply blur to image pixels
func blurImagePixels(pixels [][]Pixel) ([][]Pixel, error) {
	// Copy matrix
	blurredPixels := pixels

	for i := range pixels {
		for j := range pixels[i] {
			blurredPixels[i][j] = processBlurMask(i, j, pixels)
		}
	}

	return blurredPixels, nil
}

// Write new file with given matrix of pixels
func writeNewImage(pixels [][]Pixel, outFileName string) {
	file, err := os.Create(outFileName)
	if err != nil {
		fmt.Println("Error writing file", err)
		os.Exit(1)
	}

	width := len(pixels)
	height := len(pixels[0])

	newImage := image.NewNRGBA(image.Rect(0, 0, width, height))
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			newImage.Set(i, j, color.RGBA{uint8(pixels[i][j].R), uint8(pixels[i][j].G), uint8(pixels[i][j].B), uint8(pixels[i][j].A)})
		}
	}

	writeErr := png.Encode(file, newImage)
	if writeErr != nil {
		fmt.Println(writeErr)
		os.Exit(1)
	}
}

// Process a new image, applying the blur filter
func processImage(file io.Reader, outFileName string) {
	fmt.Println("Start processing image..")
	// Get image pixels as a matrix of Pixel
	pixelMatrix, err := getImagePixels(file)
	if err != nil {
		fmt.Println("There was an error", err)
		os.Exit(1)
	}

	fmt.Println("Applying blur filter..")
	// Generate the a blurred pixel matrix
	newMatrix, err := blurImagePixels(pixelMatrix)
	if err != nil {
		fmt.Println("Error processing image", err)
		os.Exit(1)
	}

	fmt.Println("Writing new image file..")

	// Write file
	writeNewImage(newMatrix, outFileName)

	// Finish
	fmt.Println("Done.")
	os.Exit(1)
}

// Main
func main() {
	fileName := "lenna.png"
	if len(os.Args) > 1 && string(os.Args[1]) != "" {
		fileName = os.Args[1]
	}

	outFileName := "lenna-blurred.png"
	if len(os.Args) > 2 && string(os.Args[2]) != "" {
		outFileName = string(os.Args[2])
	}

	// Read png
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	file, err := os.Open(fileName)

	if err != nil {
		fmt.Println("Error opening file: %s", fileName)
		os.Exit(1)
	}

	// Close on end
	defer file.Close()

	// Process
	processImage(file, outFileName)
}
