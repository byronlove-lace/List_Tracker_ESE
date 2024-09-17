package main

import (
	"bufio"
	"fmt"
	// "github.com/charmbracelet/bubbletea"
	"golang.org/x/net/html"
	"net/http"
	"os"
)

func main() {
	resp, err := http.Get("https://warrocketwiki.com/Every_Story_Ever")

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		respOutput := "Response Status: " + resp.Status
		panic(respOutput)
	}

	file, err := os.OpenFile("EveryStoryEverList.md", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)

	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	htmlContent := html.NewTokenizer(resp.Body)
	orderedListCount := 0
	inListItem := false
	inOrderedList := false
	text := ""

	for {
		currentToken := htmlContent.Next()
		tokenContent := htmlContent.Token()
		isOrderedList := tokenContent.Data == "ol"
		isListItem := tokenContent.Data == "li"

		switch {
		case currentToken == html.ErrorToken:
			err = writer.Flush()
			if err != nil {
				panic(err)
			}
			return

		case currentToken == html.StartTagToken:
			if isOrderedList {
				orderedListCount += 1
				inOrderedList = true
			}

			if isListItem {
				inListItem = true
				if orderedListCount == 2 {
					text = fmt.Sprintf("[ ]")
				}
			}

		case currentToken == html.TextToken:
			if inOrderedList && inListItem && orderedListCount == 2 {
				text += tokenContent.Data
			}

		case currentToken == html.EndTagToken:

			if isListItem {
				inListItem = false
				if orderedListCount == 2 && inOrderedList {
					_, err = writer.WriteString(text + "\n")
					if err != nil {
						panic(err)
					}
				}
			}

			if isOrderedList {
				inOrderedList = false
			}
		}
	}
}
