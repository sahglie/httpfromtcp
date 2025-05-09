package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":42069")
	if err != nil {
		log.Fatalf("%s\n", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("%s\n", err)
	}

	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	// scanr := bufio.NewScanner(reader)

	for {
		fmt.Print("> ")
		data, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}

		_, err = conn.Write([]byte(data))
		if err != nil {
			fmt.Println(err)
		}
	}
}
