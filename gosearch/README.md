# Web GoSearch 

## Overview
GoSearch is a command-line tool written in Go (Golang) that helps in discovering directories and analyzing HTML forms, query parameters, and JavaScript code on web pages. It is designed to assist in web application security testing by identifying potential areas of vulnerability.

## Features
- Directory enumeration: Scans a target URL with a provided wordlist to discover existing directories.
- HTML form analysis: Identifies HTML forms on web pages and analyzes their attributes, such as action URL, method, content type, and input fields.
- Query parameter analysis: Detects query parameters in URLs and provides insights for further analysis.
- JavaScript code analysis: Identifies JavaScript code embedded in web pages and assists in understanding its interaction with user input.
## How to Compile and Run
1. Clone the repository: git clone https://github.com/rosepwns/GoWasp.git
2. Navigate to the project directory: ```cd gosearch ```
3. Compile the source code: ```go build gosearch.go```
4. Run the executable with the target URL and wordlist file as arguments:
```php
# Copy code
./gosearch <url> <wordlist>
```
Replace <url> with the target URL to scan and <wordlist> with the path to the wordlist file containing directory names to test.

## Options
- -verbose: Run in verbose mode to display each directory being tested in real-time.

Example Usage
```bash
# Copy code
./gosearch https://example.com /path/to/wordlist.txt
```
## Dependencies
- Go (Golang)
- github.com/PuerkitoBio/goquery

## Contributing
Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request.
