package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// RateLimitedReader limits the read speed
// It limits how fast data can be read from the internet
type RateLimitedReader struct {
	reader      io.Reader
	bytesPerSec int64
}

// Read implements the io.Reader interface with rate limiting
/*Read Method: This method is part of the RateLimitedReader struct.
It reads data (p []byte) from the internet (saveinfo.Body).
 It calculates how fast data is being read (elapsed),
compares it to the desired speed (expectedTime),
and slows down if necessary (time.Sleep) to maintain the speed limit (bytesPerSec).*/
func (r *RateLimitedReader) Read(p []byte) (n int, err error) {
	start := time.Now()
	n, err = r.reader.Read(p)
	if n > 0 {
		elapsed := time.Since(start)
		expectedTime := time.Duration(int64(n) * int64(time.Second) / r.bytesPerSec)
		if elapsed < expectedTime {
			time.Sleep(expectedTime - elapsed)
		}
	}
	return
}

/*
parseRateLimit Function:
Converts a string like "1m" (1 megabyte per second)

	into a number of bytes per second that the download
	speed should not exceed.
*/
func parseRateLimit(rateStr string) (int64, error) {
	if len(rateStr) < 2 {
		return 0, fmt.Errorf("invalid rate limit format")
	}
	multiplier := int64(1)
	switch rateStr[len(rateStr)-1] {
	case 'k', 'K':
		multiplier = 1024
	case 'm', 'M':
		multiplier = 1024 * 1024
	default:
		return 0, fmt.Errorf("invalid rate limit suffix")
	}

	rate, err := strconv.ParseInt(rateStr[:len(rateStr)-1], 10, 64)
	if err != nil {
		return 0, err
	}

	return rate * multiplier, nil
}

/*
It uses the RateLimitedReader to control
how fast data is read and saved (io.Copy).
Finally, it returns the number of bytes downloaded.
*/
func DownloadWithSpeedLimit(fullPath string, url3 string, rate int64) (bytes int64) {
	makefile, err := os.Create(fullPath)
	Check(err)
	defer makefile.Close()

	saveinfo, err := http.Get(url3)
	Check(err)
	defer saveinfo.Body.Close()

	if saveinfo.StatusCode != 200 {
		err = fmt.Errorf("status code: %d", saveinfo.StatusCode)
		Check(err)
	}

	// Progress bar setup
	const (
		unit      = 1024
		maxLength = 50
	)
	totalSize := saveinfo.ContentLength
	startTime := time.Now()
	progress := newProgressBar(totalSize, maxLength)
	defer fmt.Println()

	limitedReader := &RateLimitedReader{
		reader:      saveinfo.Body,
		bytesPerSec: rate,
	}
	bytes, err = io.Copy(io.MultiWriter(makefile, progress), limitedReader)
	Check(err)

	elapsed := time.Since(startTime)
	speed := float64(bytes) / elapsed.Seconds()
	fmt.Printf("\nDownload completed: %s (%.2f MiB/s)\n", fullPath, speed/(unit*unit))

	return bytes
}

// Function to create a progress bar
func newProgressBar(total int64, length int) io.Writer {
	if total <= 0 {
		return os.Stdout
	}
	return &progressBar{
		total:   total,
		current: 0,
		length:  length,
		start:   time.Now(),
	}
}

// progressBar implements the io.Writer interface to show a progress bar
type progressBar struct {
	total   int64
	current int64
	length  int
	start   time.Time
}

func (p *progressBar) Write(b []byte) (int, error) {
	n := len(b)
	p.current += int64(n)

	// Calculate percentage completed
	percent := float64(p.current) / float64(p.total) * 100

	// Calculate progress bar length
	progressLength := int(float64(p.length) * (percent / 100))

	// Print progress bar
	fmt.Printf("\r[%s%s] %.2f%%",
		strings.Repeat("=", progressLength),
		strings.Repeat(" ", p.length-progressLength),
		percent)

	// Check if download is complete
	if p.current >= p.total {
		elapsed := time.Since(p.start)
		fmt.Printf(" %.2f MiB/s %s elapsed\n", float64(p.total)/elapsed.Seconds()/(1024*1024), elapsed.Truncate(time.Second))
	}

	return n, nil
}
