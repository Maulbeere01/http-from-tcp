package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const port = ":42069"

func main() {
	l, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("could not open network l on port: %s error: %s\n", port, err.Error())
	}
	defer l.Close()

	fmt.Println("listening for TCP traffic on", port)
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalf("error: %s\n", err.Error())
		}
		fmt.Printf("connection: %s accepted\n", conn.RemoteAddr())
		ch := getLinesChannel(conn)
		for line := range ch {
			fmt.Printf("read: %s\n", line)
		}
		fmt.Println("connection: ", conn.RemoteAddr(), "closed")
	}
}

func getLinesChannel(s io.ReadCloser) <-chan string {
	stringChannel := make(chan string)

	go func() {
		defer s.Close()
		defer close(stringChannel)

		buffer := make([]byte, 8)
		var sb strings.Builder

		for {
			n, err := s.Read(buffer)

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
