package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

func main() {
	// Ensure the program is run with at least one argument (section name)
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <section_name>\n", os.Args[0])
	}

	sectionName := os.Args[1]
	beginText := "### " + sectionName
	outputFileName := "observations.md"

	startTime := time.Now()

	// Define the directory path containing the content files
	dirPath := path.Join("..", "this-week-in-rust", "content")
	contentDir, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatalf("Error reading directory: %v\n", err)
	}

	var b strings.Builder

	// Iterate over each file in the directory
	for _, e := range contentDir {
		// Only process files that contain "this-week-in-rust" in their name
		if strings.Contains(e.Name(), "this-week-in-rust") {
			log.Printf("Processing file: %s\n", e.Name())
			heading := fmt.Sprintf("%s - %s", beginText, e.Name())
			data, err := os.ReadFile(path.Join(dirPath, e.Name()))
			if err != nil {
				log.Fatalf("Error reading file %s: %v\n", e.Name(), err)
			}
			content := string(data)

			// Check if the file contains the begin text
			if strings.Contains(content, beginText) {
				startOffset := strings.Index(content, beginText)
				endOffset := len(content)

				log.Printf("Found section '%s' in file: %s at position %d\n", sectionName, e.Name(), startOffset)

				// Find the start of the next section header to determine the end of the current section
				nextSectionStart := content[startOffset+len(beginText):]
				if nextSectionOffset := strings.Index(nextSectionStart, "### "); nextSectionOffset != -1 {
					endOffset = startOffset + len(beginText) + nextSectionOffset
					log.Printf("Next section starts at position %d\n", endOffset)
				}

				// Read the section content from the file
				reader := strings.NewReader(content)
				sr := io.NewSectionReader(reader, int64(startOffset), int64(endOffset)-int64(startOffset))
				buf := make([]byte, endOffset-startOffset)
				if _, err := sr.Read(buf); err != nil && err != io.EOF {
					log.Fatalf("Error reading section from file %s: %v\n", e.Name(), err)
				}

				// Replace the begin text with the heading and write it to the string builder
				newStr := strings.Replace(string(buf), beginText, heading, 1)
				if _, err := b.Write([]byte(newStr)); err != nil {
					log.Fatalf("Error writing to string builder: %v\n", err)
				}
				log.Printf("Successfully processed section from file: %s\n", e.Name())
			} else {
				log.Printf("Section '%s' not found in file: %s\n", sectionName, e.Name())
			}
		} else {
			log.Printf("Skipping file: %s\n", e.Name())
		}
	}

	// Write the collected sections to the output file
	err = os.WriteFile(outputFileName, []byte(b.String()), 0666)
	if err != nil {
		log.Fatalf("Error writing to output file %s: %v\n", outputFileName, err)
	}

	fmt.Printf("Total Time: %s\n", time.Since(startTime))
}
