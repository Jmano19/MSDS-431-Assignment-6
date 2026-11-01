# Go Image Processing Pipeline: Sequential vs. Concurrent Execution

## Overview

This project implements an image-processing pipeline in Go and compares sequential execution with a concurrent pipeline built using goroutines and channels.

The application processes four JPEG images through the following stages:

1. Read the input image.
2. Resize the image to 500 × 500 pixels.
3. Convert the image to grayscale.
4. Save the processed image as a JPEG.

The goal is to evaluate whether Go concurrency improves image-processing throughput compared with sequential execution.

## Repository Structure

```text
go_21_goroutines_pipeline/
├── image_processing/
│   ├── image_processing.go
│   └── image_processing_test.go
├── images/
│   ├── image1.jpeg
│   ├── image2.jpeg
│   ├── image3.jpeg
│   ├── image4.jpeg
│   └── output/
│       ├── sequential/
│       └── concurrent/
├── main.go
├── main_test.go
├── go.mod
├── go.sum
└── README.md
```

## Requirements

- Go 1.20 or later
- JPEG input images
- The Go dependency listed in `go.mod`

Install dependencies:

```powershell
go mod download
```

## Running the Application

The application supports three execution modes.

### Run both implementations

```powershell
go run . -mode both -runs 5
```

### Run sequential processing

```powershell
go run . -mode sequential -runs 5
```

### Run concurrent processing

```powershell
go run . -mode concurrent -runs 5
```

The `-runs` parameter specifies how many times the four-image workload is processed.

The application reports successful and failed images, elapsed time, and throughput in images per second.

Processed images are saved in separate directories:

- Sequential: `images/output/sequential/`
- Concurrent: `images/output/concurrent/`

Generated output images are excluded from version control through `.gitignore`.

## Implementation Design

### Sequential Processing

The sequential implementation processes each image from beginning to end before moving to the next image.

For each image, the application reads the file, resizes it, converts it to grayscale, and saves the result.

This approach provides a straightforward baseline without overlapping image-processing stages.

### Concurrent Processing

The concurrent implementation uses goroutines and channels to create a pipeline consisting of four stages:

1. **Load:** Reads the input images and sends jobs into a channel.
2. **Resize:** Resizes images received from the loading stage.
3. **Grayscale:** Converts resized images to grayscale.
4. **Save:** Writes the processed images to disk.

Each stage runs in its own goroutine and communicates with the next stage through channels. This allows different images to move through different stages simultaneously.

The implementation uses a staged pipeline rather than a worker pool. Each stage processes incoming jobs while the other stages can work on different images.

Errors are carried through the pipeline so that a failed image does not terminate processing of the remaining images.

### Error Handling

The application includes error handling for:

- Missing or inaccessible input files.
- Invalid or unsupported image data.
- Output directory creation failures.
- Output file creation and JPEG encoding failures.
- Nil images and invalid paths.

Errors are reported to the console, and the application maintains a summary of successful and failed images.

## Testing

Run all unit and pipeline tests:

```powershell
go test ./...
```

Run tests with coverage:

```powershell
go test -cover ./...
```

Run tests with verbose output:

```powershell
go test -v ./...
```

The tests cover:

**Image-processing package:**

- Reading valid JPEG files.
- Handling missing input files.
- Writing valid JPEG files.
- Creating output directories.
- Rejecting nil images.
- Grayscale conversion and image dimensions.
- Image resizing to 500 × 500 pixels.

**Main package:**

- Successful sequential image processing.
- Successful concurrent image processing.
- Missing-input handling in both execution modes.
- Benchmark input availability.
- Expected benchmark image count.

All 14 tests pass.

Package statement coverage:

| Package | Coverage |
|---|---:|
| Main package | 56.5% |
| Image-processing package | 84.6% |

## Performance Benchmarking

The project includes Go benchmark functions for sequential and concurrent processing.

Run the benchmarks:

```powershell
go test -run=^$ -bench='Benchmark(Sequential|Concurrent)' -benchmem -count=5
```

Each benchmark operation processes all four input images. The benchmark reports execution time, memory allocations, and custom throughput in images per second.

The results below represent the average of five benchmark repetitions.

### Benchmark Environment

| Component | Configuration |
|---|---|
| Operating System | Windows |
| Architecture | amd64 |
| Processor | AMD Ryzen 7 5800X 8-Core Processor |
| Go Version | 1.27.1 |
| Input Images | 4 JPEG files |

### Go Benchmark Results

| Metric | Sequential | Concurrent |
|---|---:|---:|
| Average time per operation | 186.50 ms | 147.47 ms |
| Average throughput | 21.45 images/sec | 27.13 images/sec |
| Memory allocated per operation | ~67.93 MB | ~67.93 MB |
| Allocations per operation | ~1,000,451 | ~1,000,444 |

The concurrent pipeline achieved approximately **26.5% higher throughput**, with an average speedup of approximately **1.26×** compared with sequential processing.

### Application-Level Results

An additional performance comparison processed 100 images per mode over five trials.

| Metric | Sequential Average | Concurrent Average |
|---|---:|---:|
| Elapsed time | 4.763 seconds | 3.647 seconds |
| Throughput | 21.00 images/sec | 27.42 images/sec |
| Successful images | 100 | 100 |
| Failed images | 0 | 0 |

This comparison showed approximately **30.5% higher throughput** with the concurrent pipeline.

### Performance Analysis

Both benchmarking approaches demonstrated that the concurrent implementation outperformed sequential processing on the tested workload.

The results suggest that overlapping image-processing stages using goroutines and channels improves throughput. Memory allocations remained nearly identical between the two implementations, indicating that the performance improvement primarily came from overlapping pipeline stages rather than reduced memory usage.

The results are specific to the tested hardware, image set, and workload. Larger image sets and additional benchmark trials could provide further insight into scalability.

## GenAI Tools

OpenAI ChatGPT (GPT-6) was used as a generative AI development assistant throughout this assignment. ChatGPT assisted with the following activities:

- **Code understanding:** Reviewed the original CODEHEIM Go image-processing pipeline and explained its use of goroutines and channels.
- **Implementation:** Assisted with restructuring the application into sequential and concurrent processing modes and implementing the staged goroutine pipeline.
- **Error handling:** Helped identify and implement checks for missing input files, invalid image data, output directory creation, and JPEG writing failures.
- **Testing:** Assisted in developing unit tests for image reading, writing, resizing, grayscale conversion, sequential processing, concurrent processing, and missing-input scenarios.
- **Benchmarking:** Helped develop Go benchmark functions, interpret benchmark output, calculate throughput improvements, and compare sequential and concurrent performance.
- **Debugging and documentation:** Assisted with troubleshooting Go compilation and execution issues and preparing the README documentation.

All AI-assisted code was reviewed, modified where necessary, and executed locally. Test results, application behavior, and benchmark measurements were verified in the local Go environment.

### GenAI Documentation

A supporting development activity log is included in [`GENAI_LOG.md`](GENAI_LOG.md), which summarizes AI-assisted development activities, and - [GenAI Conversation Log](GENAI_CONVERSATION_LOG.md), which contains representative conversation excerpts documenting the use of ChatGPT during implementation, testing, benchmarking, and documentation.

The log is a summary of the development process rather than a verbatim transcript of the conversations.

All generated code was reviewed, executed, and tested locally. The final implementation and benchmark results were verified in the local Go environment.

## Original Project

This assignment builds upon the original Go concurrency example:

[CODEHEIM – Go Goroutines Pipeline](https://github.com/code-heim/go_21_goroutines_pipeline)

## License

Refer to the original repository's LICENSE file.