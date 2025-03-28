package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const newline = 10

func main() {
	listen, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("%s\n", err)
	}

	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatalf("%s\n", err)
		}

		fmt.Println("Connection accepted")

		ch := getLinesChannel(conn)

		for line := range ch {
			fmt.Println(line)
		}

		fmt.Println("Connection closed")
	}

}

func getLinesChannel(fd io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		defer close(ch)
		defer fd.Close()

		var sb strings.Builder
		buf := make([]byte, 8, 8)

		for {
			n, err := fd.Read(buf)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}

				fmt.Printf("error: %s\n", err)
				break
			}

			for i := 0; i < n; i++ {
				if buf[i] == newline {
					ch <- sb.String()
					sb.Reset()
				} else {
					sb.WriteByte(buf[i])
				}
			}
		}

		if sb.Len() > 0 {
			ch <- sb.String()
		}
	}()

	return ch
}
