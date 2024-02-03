package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.DialUDP("udp",
		&net.UDPAddr{
			IP:   net.ParseIP("0.0.0.0"),
			Port: 0,
		},
		&net.UDPAddr{
			IP:   net.ParseIP("0.0.0.0"),
			Port: 3228,
		})
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter a message to send: ")
		if !scanner.Scan() {
			break
		}
		msg := scanner.Text()
		n, err := conn.Write([]byte(msg))
		if err != nil {
			panic(err)
		}
		fmt.Printf("%d bytes are sent\n", n)
	}

	if scanner.Err() != nil {
		panic(scanner.Err())
	}
}
