# Simplified Wget Tool

This project is a simplified implementation of the `wget` utility, designed to download files from the web. While it doesn't include all the advanced functionalities of GNU Wget, it demonstrates essential features for fetching files using HTTP(S) requests.

## Features Implemented

1. **Download a File**  
   Download a file using a given URL and save it to the current directory.

2. **Save with a Custom Filename**  
   Specify a different name for the downloaded file using the `-O` flag.

3. **Save to a Custom Directory**  
   Use the `-P` flag to specify the target directory for the downloaded file.

4. **Basic Progress Feedback**  
   Display basic download progress, including:
   - Start time
   - Status of the request
   - File size in bytes
   - Time taken to complete the download

## Usage

### Basic File Download
```bash
$ go run . https://example.com/file.txt
```

### Download and Save with a Custom Filename
```bash
$ go run . -O=custom_name.txt https://example.com/file.txt
```

### Download to a Specific Directory
```bash
$ go run . -P=/path/to/directory https://example.com/file.txt
```

### Combined Example
```bash
$ go run . -P=/downloads/ -O=custom_name.txt https://example.com/file.txt
```

### Background Download
Redirect logs to `wget-log` when using the `-B` flag:
```bash
$ go run . -B https://example.com/file.txt
$ cat wget-log
```

## Notes

- The project focuses on basic features and does not include the full suite of options provided by GNU Wget.
- Download speed limitation, asynchronous downloads, and website mirroring were not implemented in this version.

## Learning Outcomes

This project was a hands-on introduction to:
- Making HTTP(S) requests in Go
- Handling files and directories programmatically
- Providing user-friendly feedback via the command-line interface
- Building a foundational understanding of web download tools

## Credits

This project was developed in collaboration with:
- me: **Tala Amm**
- **Moaz Razem**
- **Amro Khweis**
- **Noor Halabi**
