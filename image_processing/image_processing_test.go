package imageprocessing

import (
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
	"testing"
)

// createTestImage creates a small image for unit tests.
func createTestImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 20, 10))

	for y := 0; y < 10; y++ {
		for x := 0; x < 20; x++ {
			img.Set(x, y, color.RGBA{
				R: uint8(x * 10),
				G: uint8(y * 20),
				B: 100,
				A: 255,
			})
		}
	}

	return img
}

func TestReadImage(t *testing.T) {
	tempDir := t.TempDir()
	imagePath := filepath.Join(tempDir, "test.jpg")

	file, err := os.Create(imagePath)
	if err != nil {
		t.Fatal(err)
	}

	if err := jpeg.Encode(file, createTestImage(), nil); err != nil {
		file.Close()
		t.Fatal(err)
	}
	file.Close()

	img, err := ReadImage(imagePath)
	if err != nil {
		t.Fatalf("ReadImage returned an unexpected error: %v", err)
	}

	if img == nil {
		t.Fatal("ReadImage returned a nil image")
	}

	if img.Bounds().Dx() != 20 || img.Bounds().Dy() != 10 {
		t.Errorf("unexpected image dimensions: %v", img.Bounds())
	}
}

func TestReadImageMissingFile(t *testing.T) {
	_, err := ReadImage(filepath.Join(t.TempDir(), "missing.jpg"))

	if err == nil {
		t.Fatal("expected an error when reading a missing file")
	}
}

func TestWriteImage(t *testing.T) {
	outputPath := filepath.Join(t.TempDir(), "nested", "output.jpg")

	err := WriteImage(outputPath, createTestImage())
	if err != nil {
		t.Fatalf("WriteImage returned an unexpected error: %v", err)
	}

	file, err := os.Open(outputPath)
	if err != nil {
		t.Fatalf("output file was not created: %v", err)
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		t.Fatalf("output file is not a valid JPEG: %v", err)
	}

	if img.Bounds().Dx() != 20 || img.Bounds().Dy() != 10 {
		t.Errorf("unexpected output dimensions: %v", img.Bounds())
	}
}

func TestWriteImageNil(t *testing.T) {
	err := WriteImage(filepath.Join(t.TempDir(), "output.jpg"), nil)

	if err == nil {
		t.Fatal("expected an error when writing a nil image")
	}
}

func TestGrayscale(t *testing.T) {
	img := createTestImage()
	grayImg := Grayscale(img)

	if grayImg == nil {
		t.Fatal("Grayscale returned a nil image")
	}

	if grayImg.Bounds() != img.Bounds() {
		t.Errorf("expected bounds %v, got %v", img.Bounds(), grayImg.Bounds())
	}

	r, g, b, _ := grayImg.At(5, 5).RGBA()

	if r != g || g != b {
		t.Errorf("pixel is not grayscale: R=%d G=%d B=%d", r, g, b)
	}
}

func TestGrayscaleNil(t *testing.T) {
	if Grayscale(nil) != nil {
		t.Error("expected Grayscale(nil) to return nil")
	}
}

func TestResize(t *testing.T) {
	img := createTestImage()
	resized := Resize(img)

	if resized == nil {
		t.Fatal("Resize returned a nil image")
	}

	if resized.Bounds().Dx() != 500 || resized.Bounds().Dy() != 500 {
		t.Errorf("expected 500x500 image, got %dx%d",
			resized.Bounds().Dx(), resized.Bounds().Dy())
	}
}

func TestResizeNil(t *testing.T) {
	if Resize(nil) != nil {
		t.Error("expected Resize(nil) to return nil")
	}
}
