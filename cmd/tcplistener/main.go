package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"log"
	"net"
)

//import (
//	"fmt"
//	"httpfromtcp/internal/request"
//	"log"
//	"net"
//)

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

		req, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Printf("error: %v\n", err)
		} else {
			rl := req.RequestLine
			fmt.Println("Request line:")
			fmt.Printf("- Method: %s\n", rl.Method)
			fmt.Printf("- Target: %s\n", rl.RequestTarget)
			fmt.Printf("- Version: %s\n", rl.HttpVersion)
			headers := req.Headers
			fmt.Println("Headers:")
			for k, v := range headers {
				fmt.Printf("- %s: %s\n", k, v)
			}
			fmt.Println("Body:")
			fmt.Printf("%s\n", req.Body)
		}

		fmt.Println("Connection closed")
	}
}
