package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	gg string = "./"
	// myurl  string
	mybyte int64
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Please provide a url to download")
		return
	}
	allArgs := os.Args[1:]
	myurl := os.Args[len(os.Args)-1]
	if bol("--mirror", allArgs) == true {
		mirrorTheFuckingWEBISATE(myurl)
	}

	filename := FindFileName(myurl)

	flagO := flag.String("O", filename, "Usage -O=Filename")
	flagP := flag.String("P", "./", "Usage -P=FilePath")
	flagB := flag.Bool("B", false, "Usage -b")
	flagI := flag.String("i", "", "Usage: -i=links.txt")
	flagRate := flag.String("rate-limit", "", "Rate limit for download (e.g., 200k, 2M)")
	flag.Parse()

	old := os.Stdout
	defer func() { os.Stdout = old }() // Restore original stdout

	if *flagB == true {
		fmt.Println("output will be written in ./wget-log")
		t, err := os.Create("wget-log")
		Check(err)
		defer t.Close()
		os.Stdout = t
	}

	if *flagI != "" {

		read, err := os.ReadFile(*flagI)
		Check(err)
		stringRead := string(read)
		sp := strings.Split(stringRead, "\n")

		for _, oneURL := range sp {
			start1 := time.Now()
			fmt.Println("Starting download at", start1.Format("2006-01-02 15:04:05"))
			oneURL = strings.TrimSpace(oneURL)
			// fmt.Println([]byte(oneURL))
			namefile := FindFileName(oneURL)
			mybyte, err = DownloadLink(gg+namefile, oneURL)
			Check(err)
			end1 := time.Now()

			fmt.Println("Finshed download at", end1.Format("2006-01-02 15:04:05"))

			totalTime := end1.Sub(start1)
			fmt.Println("DowmLoaded url:", oneURL)
			fmt.Println("sending request took", totalTime)
			fmt.Printf("content size %d [~%.2fMB]\n", mybyte, float64(mybyte)/(1024*1024))

			// fmt.Println(PrintTime(start1, end1, oneURL))
		}
		return
	}

	if *flagP != "./" {
		home, err := os.UserHomeDir() // C:\Users\lenovo     C:\Users\adam
		Check(err)
		j := *flagP
		gg = filepath.Join(home, j[1:]) // -P=~/Downloads/   -->   C:\Users\adam\Downloads\
	}

	merge := filepath.Join(gg, *flagO) // connects the path with the file name --> C:\Users\adam\Downloads\EMtmPFLWkAA8CIS.jpg
	if *flagRate != "" {
		speedRate, err := parseRateLimit(*flagRate)
		Check(err)
		start := time.Now()
		fmt.Println("Starting download at", start.Format("2006-01-02 15:04:05"))
		mybyte = DownloadWithSpeedLimit(merge, myurl, speedRate)
		end := time.Now()
		fmt.Println("Finshed download at", end.Format("2006-01-02 15:04:05"))

		totalTime := end.Sub(start)
		fmt.Println("DowmLoaded url:", myurl)
		fmt.Println("sending request took", totalTime)
		fmt.Printf("content size %d [~%.2fMB]\n", mybyte, float64(mybyte)/(1024*1024))

		//	fmt.Println(PrintTime(start, end, myurl))

		os.Exit(0)
	}

	startTime := time.Now()
	fmt.Println("Starting download at", startTime.Format("2006-01-02 15:04:05"))
	mybyte, err := DownloadLink(merge, myurl)
	Check(err)
	endtime := time.Now()
	fmt.Println("Finshed download at", endtime.Format("2006-01-02 15:04:05"))

	totalTime := endtime.Sub(startTime)
	fmt.Println("DowmLoaded url:", myurl)
	fmt.Println("sending request took", totalTime)
	fmt.Printf("content size %d [~%.2fMB]\n", mybyte, float64(mybyte)/(1024*1024))
}

func DownloadLink(fullPath string, url3 string) (bytes int64, err error) {
	// Create file to save the downloaded content
	makefile, err := os.Create(fullPath)
	if err != nil {
		return 0, err
	}
	defer makefile.Close()

	// Make HTTP GET request to download the content
	resp, err := http.Get(url3)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("HTTP status code: %d", resp.StatusCode)
	}

	// Create a progress bar
	fmt.Printf("Downloading %s...\n", fullPath)
	const (
		unit      = 1024
		maxLength = 50 // Length of the progress bar
	)
	totalSize := resp.ContentLength
	startTime := time.Now()
	progress := newProgressBar(totalSize, maxLength)
	defer fmt.Println() // Print a new line after the progress bar completes

	// Copy content from HTTP response to the file with progress tracking
	bytesCopied, err := io.Copy(io.MultiWriter(makefile, progress), resp.Body)
	if err != nil {
		return 0, err
	}

	// Calculate and print download statistics
	elapsed := time.Since(startTime)
	speed := float64(bytesCopied) / elapsed.Seconds()
	fmt.Printf("\nDownload completed: %s (%.2f MiB/s)\n", fullPath, speed/(unit*unit))

	return bytesCopied, nil
}

func FindFileName(url1 string) string {
	split := strings.Split(url1, "/")
	filename := split[len(split)-1] // notice in the above line, the name of the file is the last element in the slice split
	return filename
}

func Check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

func bol(s string, arr []string) bool {
	for _, k := range arr {
		if s == k {
			return true
		}
	}
	return false
}
