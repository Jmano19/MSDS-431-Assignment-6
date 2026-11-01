package imageprocessing

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
)

// ReadImage opens and decodes an image from the provided file path.
// It returns an error instead of terminating the application.
func ReadImage(path string) (image.Image, error) {
	if path == "" {
		return nil, fmt.Errorf("image path cannot be empty")
	}

	inputFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open image %q: %w", path, err)
	}
	defer inputFile.Close()

	img, _, err := image.Decode(inputFile)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image %q: %w", path, err)
	}

	return img, nil
}

// WriteImage creates the output directory if necessary
// and saves the image as a JPEG file.
func WriteImage(path string, img image.Image) error {
	if path == "" {
		return fmt.Errorf("output path cannot be empty")
	}

	if img == nil {
		return fmt.Errorf("cannot write a nil image")
	}

	outputDir := filepath.Dir(path)

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	outputFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create output file %q: %w", path, err)
	}

	if err := jpeg.Encode(outputFile, img, nil); err != nil {
		outputFile.Close()
		return fmt.Errorf("failed to encode image %q: %w", path, err)
	}

	if err := outputFile.Close(); err != nil {
		return fmt.Errorf("failed to close output file %q: %w", path, err)
	}

	return nil
}

// Grayscale converts an image into grayscale.
func Grayscale(img image.Image) image.Image {
	if img == nil {
		return nil
	}

	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalPixel := img.At(x, y)
			grayPixel := color.GrayModel.Convert(originalPixel)
			grayImg.Set(x, y, grayPixel)
		}
	}

	return grayImg
}

// Resize resizes an image to 500 x 500 pixels.
// Aspect ratio preservation will be implemented separately.
func Resize(img image.Image) image.Image {
	if img == nil {
		return nil
	}

	newWidth := uint(500)
	newHeight := uint(500)

	return resize.Resize(newWidth, newHeight, img, resize.Lanczos3)
}
