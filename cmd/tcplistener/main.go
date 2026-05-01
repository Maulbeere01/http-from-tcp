package main

import (
	"fmt"
	"log"
	"net"

	"httpfromtcp/internal/request"
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
		r, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("error: %s\n", err.Error())
		}
		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\nHeaders:\n",
			r.RequestLine.Method, r.RequestLine.RequestTarget, r.RequestLine.HttpVersion)
		for key, value := range r.Headers {
			fmt.Printf("- %s: %s\n", key, value)
		}
	}
}
