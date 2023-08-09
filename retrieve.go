package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"strings"
)

func main() {
	// URL of the USCF profile
	url := "http://www.uschess.org/msa/MbrDtlMain.php?15518872"

	// Send a GET request to the URL
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer response.Body.Close()

	// Check if the request was successful
	if response.StatusCode == 200 {
		// Parse the HTML content using the HTML package
		tokenizer := html.NewTokenizer(response.Body)

		var playerName string
		var rating string
		var inPlayerName, inRating bool

		for {
			tokenType := tokenizer.Next()
			switch tokenType {
			case html.ErrorToken:
				if tokenizer.Err() == io.EOF {
					// Finished parsing
					fmt.Println("Player Name:", playerName)
					fmt.Println("Rating:", rating)
					return
				}
			case html.TextToken:
				token := tokenizer.Token()
				text := strings.TrimSpace(token.Data)

				if inPlayerName {
					playerName = text
				} else if inRating {
					rating = text
				}
			case html.StartTagToken, html.EndTagToken:
				token := tokenizer.Token()
				if token.Data == "font" {
					for _, attr := range token.Attr {
						if attr.Key == "size" && attr.Val == "+2" {
							inPlayerName = token.Type == html.StartTagToken
						}
						if attr.Key == "color" && attr.Val == "#FF6600" {
							inRating = token.Type == html.StartTagToken
						}
					}
				}
			}
		}
	} else {
		fmt.Println("Failed to retrieve the webpage")
	}
}
