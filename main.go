package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const inputFile = "messages.txt"

func main() {
	f, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("could not open %s: %s\n", inputFile, err)
	}

	ch := getLinesChannel(f)

	for line := range ch {
		fmt.Printf("read: %s\n", line)
	}
}
func getLinesChannel(f io.ReadCloser) <-chan string {
	stringChannel := make(chan string)

	go func() {
		defer f.Close()
		defer close(stringChannel)

		buffer := make([]byte, 8)
		var sb strings.Builder

		for {
			n, err := f.Read(buffer)

			lines := strings.Split(string(buffer[:n]), "\n")
			numLines := len(lines)

			sb.WriteString(lines[0])

			if numLines > 1 {
				stringChannel <- sb.String()

				for i := 1; i < numLines-1; i++ {
					stringChannel <- lines[i]
				}
				sb.Reset()
				sb.WriteString(lines[numLines-1])
			}

			if err != nil {
				if errors.Is(err, io.EOF) {
					if sb.Len() > 0 {
						stringChannel <- sb.String()
					}
				}
				break
			}
		}
	}()
	return stringChannel
}
