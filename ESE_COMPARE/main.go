package main

import (
	"bufio"
	"log"
	"os"
)

const currentFileName = "ESE.md"
const updateFileName = "update.md"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func countLines(filename string) int {
	lineCount := 0

	file, err := os.Open(filename)
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lineCount++
	}

	return lineCount
}

func readChoicesFromFile(filename string, entryCount int) []string {
	choices := make([]string, entryCount)
	currentChoice := 0

	file, err := os.Open(filename)
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		choices[currentChoice] = scanner.Text()
		currentChoice++
	}

	return choices
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	logFile, err := os.OpenFile("compare.log", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	check(err)
	defer logFile.Close()
	log.SetOutput(logFile)

	currentChoiceCount := countLines(currentFileName)
	updateChoiceCount := countLines(updateFileName)
	log.Println("currentChoiceCount:", currentChoiceCount)
	log.Println("updateChoiceCount:", updateChoiceCount)

	currentChoices := readChoicesFromFile(currentFileName, currentChoiceCount)
	updateChoices := readChoicesFromFile(updateFileName, updateChoiceCount)

	if len(updateChoices) > len(currentChoices) {
		for updateIndex, updateVal := range updateChoices {
			for _, currentVal := range currentChoices {
				if updateVal[3:] == currentVal[3:] && currentVal[:3] == "[X]" {
					log.Println("Match with selection found:", currentVal)
					updateChoices[updateIndex] = currentVal
				}
			}
		}

		file, err := os.OpenFile(currentFileName, os.O_WRONLY|os.O_TRUNC, 0644)
		check(err)
		defer file.Close()

		writer := bufio.NewWriter(file)

		for _, line := range updateChoices {
			_, err := writer.WriteString(line + "\n")
			check(err)
		}

		writer.Flush()

		log.Println("updateChoices length:", len(updateChoices))
		log.Println("currentFile length after update:", countLines(currentFileName))
	}
}
