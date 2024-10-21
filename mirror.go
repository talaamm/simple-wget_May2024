package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

func mirrorTheFuckingWEBISATE(inputURL string) {
	resp, err := http.Get(inputURL)
	if err != nil {
		fmt.Println("Error fetching URL:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 && strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		fmt.Println("Website is mirrored")
		htmlContent, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading HTML content:", err)
			return
		}

		parsedURL, err := url.Parse(inputURL)
		if err != nil {
			fmt.Println("Error parsing URL:", err)
			return
		}
		siteDir := parsedURL.Host
		err = os.MkdirAll(siteDir, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}
		// Save the index.html file
		indexFilePath := filepath.Join(siteDir, "index.html")
		err = os.WriteFile(indexFilePath, htmlContent, 0o644)
		// err = os.WriteFile("index.html", htmlContent, 0o644)
		if err != nil {
			fmt.Println("Error saving index.html:", err)
			return
		}

		tokenizer := html.NewTokenizer(strings.NewReader(string(htmlContent)))
		for {
			tokenType := tokenizer.Next()
			if tokenType == html.ErrorToken {
				break
			}

			parsedURL, err := url.Parse(inputURL)
			if err != nil {
				fmt.Println("Error parsing URL:", err)
				return
			}
			token := tokenizer.Token()

			if token.Type == html.StartTagToken {
				for _, attr := range token.Attr {
					// if attr.Key == "src" && (strings.Contains(attr.Val, "js") || strings.Contains(attr.Val, "css")) {
					if (token.Data == "link" && attr.Key == "href" && strings.Contains(attr.Val, "css")) ||
						(token.Data == "script" && attr.Key == "src" && strings.Contains(attr.Val, "js")) {
						parsedSrc, err := url.Parse(attr.Val)
						if err != nil {
							fmt.Println("Error parsing src attribute:", err)
							continue
						}

						// Resolve relative URLs
						if !parsedSrc.IsAbs() {
							attr.Val = parsedURL.ResolveReference(parsedSrc).String()
						}

						// Create directory for assets
						fileName := filepath.Base(attr.Val)
						filePath := filepath.Join(siteDir, fileName)

						// Determine file name and path
						// fileName := filepath.Base(attr.Val)
						// filePath := filepath.Join(assetDir, fileName)

						// Download the asset
						_, err = DownloadLink(filePath, attr.Val)
						if err != nil {
							fmt.Println("Error downloading file:", err)
							continue
						}

						fmt.Printf("Downloaded: %s to %s\n", attr.Val, filePath)
					}
				}
			}
		}
	} else {
		fmt.Println("Failed to retrieve a valid HTML document")
	}
	os.Exit(0)
}
