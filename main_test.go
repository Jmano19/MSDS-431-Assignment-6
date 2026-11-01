package main

import (
	"os"
	"path/filepath"
	"testing"
)

var benchmarkImagePaths = []string{
	"images/image1.jpeg",
	"images/image2.jpeg",
	"images/image3.jpeg",
	"images/image4.jpeg",
}

// verifyBenchmarkInputs makes sure all benchmark input files exist.
func verifyBenchmarkInputs(b *testing.B) {
	b.Helper()

	for _, path := range benchmarkImagePaths {
		if _, err := os.Stat(path); err != nil {
			b.Fatalf("benchmark input %q is unavailable: %v", path, err)
		}
	}
}

// verifyOutputFiles checks that each expected output file exists.
func verifyOutputFiles(t *testing.T, paths []string, outputDir string) {
	t.Helper()

	for _, path := range paths {
		outputPath := filepath.Join(outputDir, filepath.Base(path))

		info, err := os.Stat(outputPath)
		if err != nil {
			t.Errorf("expected output file %q: %v", outputPath, err)
			continue
		}

		if info.Size() == 0 {
			t.Errorf("output file %q is empty", outputPath)
		}
	}
}

func TestProcessSequential(t *testing.T) {
	outputDir := t.TempDir()

	successes, failures := processSequential(benchmarkImagePaths, outputDir)

	if successes != len(benchmarkImagePaths) {
		t.Errorf("expected %d successes, got %d",
			len(benchmarkImagePaths), successes)
	}

	if failures != 0 {
		t.Errorf("expected 0 failures, got %d", failures)
	}

	verifyOutputFiles(t, benchmarkImagePaths, outputDir)
}

func TestProcessConcurrent(t *testing.T) {
	outputDir := t.TempDir()

	successes, failures := processConcurrent(benchmarkImagePaths, outputDir)

	if successes != len(benchmarkImagePaths) {
		t.Errorf("expected %d successes, got %d",
			len(benchmarkImagePaths), successes)
	}

	if failures != 0 {
		t.Errorf("expected 0 failures, got %d", failures)
	}

	verifyOutputFiles(t, benchmarkImagePaths, outputDir)
}

func TestProcessSequentialMissingInput(t *testing.T) {
	outputDir := t.TempDir()

	paths := []string{
		benchmarkImagePaths[0],
		filepath.Join(t.TempDir(), "missing.jpeg"),
	}

	successes, failures := processSequential(paths, outputDir)

	if successes != 1 {
		t.Errorf("expected 1 success, got %d", successes)
	}

	if failures != 1 {
		t.Errorf("expected 1 failure, got %d", failures)
	}

	verifyOutputFiles(t, paths[:1], outputDir)
}

func TestProcessConcurrentMissingInput(t *testing.T) {
	outputDir := t.TempDir()

	paths := []string{
		benchmarkImagePaths[0],
		filepath.Join(t.TempDir(), "missing.jpeg"),
	}

	successes, failures := processConcurrent(paths, outputDir)

	if successes != 1 {
		t.Errorf("expected 1 success, got %d", successes)
	}

	if failures != 1 {
		t.Errorf("expected 1 failure, got %d", failures)
	}

	verifyOutputFiles(t, paths[:1], outputDir)
}

func TestBenchmarkInputsExist(t *testing.T) {
	for _, path := range benchmarkImagePaths {
		if _, err := os.Stat(path); err != nil {
			t.Errorf("input image %q is unavailable: %v", path, err)
		}
	}
}

func TestBenchmarkImageCount(t *testing.T) {
	if len(benchmarkImagePaths) != 4 {
		t.Errorf("expected 4 benchmark images, got %d",
			len(benchmarkImagePaths))
	}
}

// BenchmarkSequential measures the sequential image-processing pipeline.
func BenchmarkSequential(b *testing.B) {
	verifyBenchmarkInputs(b)

	outputDir := b.TempDir()

	b.ReportAllocs()
	b.ReportMetric(float64(len(benchmarkImagePaths)), "images/op")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		successes, failures := processSequential(
			benchmarkImagePaths,
			outputDir,
		)

		if failures != 0 || successes != len(benchmarkImagePaths) {
			b.Fatalf(
				"sequential processing failed: %d successes, %d failures",
				successes,
				failures,
			)
		}
	}

	b.ReportMetric(
		float64(b.N*len(benchmarkImagePaths))/b.Elapsed().Seconds(),
		"images/sec",
	)
}

// BenchmarkConcurrent measures the goroutine-based pipeline.
func BenchmarkConcurrent(b *testing.B) {
	verifyBenchmarkInputs(b)

	outputDir := b.TempDir()

	b.ReportAllocs()
	b.ReportMetric(float64(len(benchmarkImagePaths)), "images/op")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		successes, failures := processConcurrent(
			benchmarkImagePaths,
			outputDir,
		)

		if failures != 0 || successes != len(benchmarkImagePaths) {
			b.Fatalf(
				"concurrent processing failed: %d successes, %d failures",
				successes,
				failures,
			)
		}
	}

	b.ReportMetric(
		float64(b.N*len(benchmarkImagePaths))/b.Elapsed().Seconds(),
		"images/sec",
	)
}
