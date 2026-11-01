# GenAI Development Log

## Project Information

- **Project:** MSDS 431 Assignment 6 – Go Image Processing Pipeline
- **Repository:**&#x20;

  [https://github.com/Jmano19/MSDS-431-Assignment-6](https://github.com/Jmano19/MSDS-431-Assignment-6)
- **AI Tool:** OpenAI ChatGPT (GPT-6)
- **Purpose:** Development assistance, code review, testing, benchmarking, debugging, and documentation.

## AI-Assisted Development Activities

### 1. Understanding the Original Project

ChatGPT assisted in reviewing the original CODEHEIM Go image-processing pipeline.

Topics discussed:

- Understanding the original image-processing functions.
- Understanding goroutines and channels.
- Identifying how the original pipeline processes images.
- Planning sequential and concurrent implementations.

### 2. Implementation and Error Handling

ChatGPT assisted with designing and reviewing changes to the Go implementation.

Topics discussed:

- Creating sequential and concurrent execution modes.
- Implementing a staged pipeline using goroutines and channels.
- Passing image-processing jobs between pipeline stages.
- Handling missing input files and invalid image data.
- Handling output directory and file-writing errors.
- Reporting successful and failed image-processing operations.

### 3. Unit Testing

ChatGPT assisted in developing tests for the image-processing package and main application.

Testing included:

- Reading valid images and handling missing files.
- Writing JPEG images and rejecting nil images.
- Grayscale conversion.
- Image resizing.
- Sequential pipeline execution.
- Concurrent pipeline execution.
- Missing-input behavior.
- Benchmark input validation.

Local verification:

```powershell
go test ./...
go test -cover ./...
go test -v ./...
```

All 14 tests passed.

Reported statement coverage:

- Main package: 56.5%.
- Image-processing package: 84.6%.

### 4. Benchmarking and Performance Analysis

ChatGPT assisted in developing benchmark functions and interpreting sequential versus concurrent performance results.

The benchmark process included:

- Measuring sequential processing time.
- Measuring concurrent pipeline processing time.
- Reporting memory allocations.
- Calculating image throughput.
- Comparing five benchmark repetitions.
- Interpreting the performance differences.

Average Go benchmark results:

| Metric             |       Sequential |       Concurrent |
| ------------------ | ---------------: | ---------------: |
| Time per operation |        186.50 ms |        147.47 ms |
| Throughput         | 21.45 images/sec | 27.13 images/sec |
| Memory allocated   |       \~67.93 MB |       \~67.93 MB |

The concurrent implementation achieved approximately 26.5% higher throughput in the Go benchmark results.

### 5. Debugging and Final Verification

ChatGPT assisted with troubleshooting and interpreting local Go command output.

The following commands were successfully executed:

```powershell
go test -v ./...
go vet ./...
go build ./...
```

The race detector was also attempted:

```powershell
go test -race ./...
```

It could not run because CGO was disabled and GCC was unavailable in the Windows environment.

### 6. Documentation and GitHub Submission

ChatGPT assisted with:

- Preparing the README.
- Documenting implementation design and benchmark results.
- Reviewing Git staging and generated output exclusions.
- Configuring the personal GitHub repository as the `origin` remote.
- Reviewing the final commit and push process.

The completed project was committed and pushed to:

[https://github.com/Jmano19/MSDS-431-Assignment-6](https://github.com/Jmano19/MSDS-431-Assignment-6)

## Verification and Responsibility

ChatGPT was used as an assistant, not as a replacement for local development and verification. Generated suggestions and code were reviewed and tested in the local Go environment. The reported test results, coverage figures, and benchmark measurements came from actual local command execution.

This document is a summary of AI-assisted development activities, not a verbatim transcript of the ChatGPT conversation history.
