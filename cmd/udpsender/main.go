package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":42069")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving UDP Address: %v", err)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error dialing UDP: %v", err)
		os.Exit(1)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		data, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v", err)
			os.Exit(1)
		}

		_, err = conn.Write([]byte(data))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error sending message: %v", err)
			os.Exit(1)
		}
	}
}
