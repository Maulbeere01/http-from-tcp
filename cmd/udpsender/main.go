package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const address = "localhost:42069"
func main(){
	addr, err := net.ResolveUDPAddr("udp",address)
	if(err != nil){
		log.Fatalf(err.Error())
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if(err != nil){
		log.Fatalf(err.Error())
	}
	defer conn.Close()

	r := bufio.NewReader(os.Stdin)
	for{
		fmt.Printf(">")
		line, err := r.ReadString('\n')
		if(err != nil){
			log.Fatalf(err.Error())
		}

		_,err = conn.Write([]byte(line))
		if(err != nil){
			log.Fatalf(err.Error())
		}
	}
}