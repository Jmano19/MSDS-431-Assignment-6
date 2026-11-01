package main

import (
	"flag"
	"fmt"
	"image"
	"path/filepath"
	"time"

	imageprocessing "goroutines_pipeline/image_processing"
)

type Job struct {
	InputPath string
	Image     image.Image
	OutPath   string
	Err       error
}

// loadImage reads the input images and sends jobs into the pipeline.
func loadImage(paths []string, outputDir string) <-chan Job {
	out := make(chan Job)

	go func() {
		defer close(out)

		for _, path := range paths {
			job := Job{
				InputPath: path,
				OutPath:   filepath.Join(outputDir, filepath.Base(path)),
			}

			job.Image, job.Err = imageprocessing.ReadImage(path)
			out <- job
		}
	}()

	return out
}

// resize processes each image in its own pipeline stage.
func resize(input <-chan Job) <-chan Job {
	out := make(chan Job)

	go func() {
		defer close(out)

		for job := range input {
			if job.Err == nil {
				job.Image = imageprocessing.Resize(job.Image)
				if job.Image == nil {
					job.Err = fmt.Errorf("resize returned a nil image: %s", job.InputPath)
				}
			}
			out <- job
		}
	}()

	return out
}

// convertToGrayscale converts each image in its own pipeline stage.
func convertToGrayscale(input <-chan Job) <-chan Job {
	out := make(chan Job)

	go func() {
		defer close(out)

		for job := range input {
			if job.Err == nil {
				job.Image = imageprocessing.Grayscale(job.Image)
				if job.Image == nil {
					job.Err = fmt.Errorf("grayscale returned a nil image: %s", job.InputPath)
				}
			}
			out <- job
		}
	}()

	return out
}

// saveImage writes completed images and sends one result per job.
func saveImage(input <-chan Job) <-chan error {
	out := make(chan error)

	go func() {
		defer close(out)

		for job := range input {
			if job.Err == nil {
				job.Err = imageprocessing.WriteImage(job.OutPath, job.Image)
			}
			out <- job.Err
		}
	}()

	return out
}

// processSequential processes each image from start to finish before
// moving to the next image.
func processSequential(paths []string, outputDir string) (int, int) {
	successCount := 0
	failureCount := 0

	for _, path := range paths {
		img, err := imageprocessing.ReadImage(path)
		if err == nil {
			img = imageprocessing.Resize(img)
			if img == nil {
				err = fmt.Errorf("resize returned a nil image: %s", path)
			}
		}

		if err == nil {
			img = imageprocessing.Grayscale(img)
			if img == nil {
				err = fmt.Errorf("grayscale returned a nil image: %s", path)
			}
		}

		if err == nil {
			outputPath := filepath.Join(outputDir, filepath.Base(path))
			err = imageprocessing.WriteImage(outputPath, img)
		}

		if err != nil {
			fmt.Printf("Sequential failure: %v\n", err)
			failureCount++
		} else {
			successCount++
		}
	}

	return successCount, failureCount
}

// processConcurrent runs the channel-based goroutine pipeline.
func processConcurrent(paths []string, outputDir string) (int, int) {
	channel1 := loadImage(paths, outputDir)
	channel2 := resize(channel1)
	channel3 := convertToGrayscale(channel2)
	writeResults := saveImage(channel3)

	successCount := 0
	failureCount := 0

	for err := range writeResults {
		if err != nil {
			fmt.Printf("Concurrent failure: %v\n", err)
			failureCount++
		} else {
			successCount++
		}
	}

	return successCount, failureCount
}

// reportResults displays elapsed time and processing throughput.
func reportResults(mode string, elapsed time.Duration, successes, failures int) {
	total := successes + failures
	seconds := elapsed.Seconds()

	fmt.Printf("\n%s Results\n", mode)
	fmt.Println("------------------------")
	fmt.Printf("Successful: %d\n", successes)
	fmt.Printf("Failed:     %d\n", failures)
	fmt.Printf("Total:      %d\n", total)
	fmt.Printf("Elapsed:    %v\n", elapsed)

	if seconds > 0 {
		fmt.Printf("Throughput: %.2f images/second\n", float64(total)/seconds)
	}
}

func main() {
	mode := flag.String("mode", "both", "Processing mode: sequential, concurrent, or both")
	runs := flag.Int("runs", 5, "Number of times to process the input images")
	flag.Parse()

	if *runs < 1 {
		fmt.Println("Error: runs must be at least 1")
		return
	}

	imagePaths := []string{
		"images/image1.jpeg",
		"images/image2.jpeg",
		"images/image3.jpeg",
		"images/image4.jpeg",
	}

	if *mode != "sequential" && *mode != "concurrent" && *mode != "both" {
		fmt.Println("Error: mode must be sequential, concurrent, or both")
		return
	}

	fmt.Printf("Processing %d images for %d run(s).\n", len(imagePaths), *runs)

	if *mode == "sequential" || *mode == "both" {
		start := time.Now()
		successes := 0
		failures := 0

		for i := 0; i < *runs; i++ {
			s, f := processSequential(imagePaths, filepath.Join("images", "output", "sequential"))
			successes += s
			failures += f
		}

		reportResults("Sequential", time.Since(start), successes, failures)
	}

	if *mode == "concurrent" || *mode == "both" {
		start := time.Now()
		successes := 0
		failures := 0

		for i := 0; i < *runs; i++ {
			s, f := processConcurrent(imagePaths, filepath.Join("images", "output", "concurrent"))
			successes += s
			failures += f
		}

		reportResults("Concurrent", time.Since(start), successes, failures)
	}
}
