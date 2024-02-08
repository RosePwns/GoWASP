package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var silent bool

// ANSI color codes
const (
	colorBlue  = "\033[1;34m"
	colorRed   = "\033[1;31m"
	colorReset = "\033[0m"
)

func main() {
	flag.BoolVar(&silent, "silent", false, "Run in silent mode (do not show directory testing messages)")
	flag.Parse()

	if len(flag.Args()) != 2 {
		fmt.Println("Usage: ./webdirfinder <url> <wordlist>")
		os.Exit(1)
	}

	baseURL := flag.Arg(0)
	wordlistFile := flag.Arg(1)

	wordlist, err := loadWordlist(wordlistFile)
	if err != nil {
		fmt.Println("Error loading wordlist:", err)
		os.Exit(1)
	}

	directories, err := getDirectories(baseURL, wordlist)
	if err != nil {
		fmt.Println("Error getting directories:", err)
		os.Exit(1)
	}

	for _, dir := range directories {
		fmt.Println(colorBlue + dir + colorReset)
		findUserInputAreas(baseURL, dir)
	}
}

func loadWordlist(filename string) ([]string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	var wordlist []string
	for _, line := range lines {
		word := strings.TrimSpace(line)
		if word != "" {
			wordlist = append(wordlist, word)
		}
	}
	return wordlist, nil
}

func getDirectories(baseURL string, wordlist []string) ([]string, error) {
	var directories []string
	for _, word := range wordlist {
		fullURL := fmt.Sprintf("%s/%s", baseURL, word)
		resp, err := http.Get(fullURL)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		// Check for final URL after following redirects
		if resp.Request.URL.String() == fullURL && resp.StatusCode == http.StatusOK {
			directories = append(directories, fullURL)
			if !silent {
				fmt.Println(colorBlue + fullURL + colorReset)
			}
		}
	}
	return directories, nil
}

func findUserInputAreas(baseURL, url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error getting page:", err)
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("Error parsing document:", err)
		return
	}

	// Look for HTML forms
	doc.Find("form").Each(func(i int, form *goquery.Selection) {
		action, exists := form.Attr("action")
		if exists {
			// Check if the form action is within the same domain as the base URL
			if strings.HasPrefix(action, baseURL) {
				fmt.Printf(colorRed+"Found form with action: %s\n"+colorReset, action)
				method := form.AttrOr("method", "GET")
				if strings.ToUpper(method) == "POST" {
					// Check for various content types
					contentType := form.AttrOr("enctype", "")
					switch strings.ToLower(contentType) {
					case "application/json":
						fmt.Println(colorRed + "Found form submitting JSON data" + colorReset)
					case "multipart/form-data":
						fmt.Println(colorRed + "Found form submitting multipart form data" + colorReset)
					case "application/x-www-form-urlencoded":
						fmt.Println(colorRed + "Found form submitting urlencoded data" + colorReset)
					case "text/plain":
						fmt.Println(colorRed + "Found form submitting plain text data" + colorReset)
					default:
						fmt.Printf(colorRed+"Found form with unknown content type: %s\n"+colorReset, contentType)
					}
				}
				// Additional logic to analyze form inputs can be added here
				form.Find("input, textarea, select").Each(func(j int, input *goquery.Selection) {
					fmt.Printf(colorRed+"Found input field: %s\n"+colorReset, input.AttrOr("name", "Unnamed"))
				})
				form.Find("button").Each(func(j int, button *goquery.Selection) {
					fmt.Printf(colorRed+"Found button: %s\n"+colorReset, button.Text())
				})
			}
		}
	})

	// Look for query parameters in URLs
	doc.Find("a").Each(func(i int, link *goquery.Selection) {
		href, exists := link.Attr("href")
		if exists {
			// Check if the link is within the same domain as the base URL
			if strings.HasPrefix(href, baseURL) {
				if strings.Contains(href, "?") {
					fmt.Printf(colorRed+"Found link with query parameters: %s\n"+colorReset, href)
					// Additional logic to analyze query parameters can be added here
				}
			}
		}
	})

	// Look for JavaScript code that interacts with user input
	doc.Find("script").Each(func(i int, script *goquery.Selection) {
		// Additional logic to analyze JavaScript code can be added here
	})
}
