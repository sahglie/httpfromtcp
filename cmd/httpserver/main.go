package main

import (
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const port = 42069

func main() {
	server, err := server.Serve(port, myHandler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Server gracefully stopped")
}

func myHandler(w *response.Writer, req *request.Request) {
	html := handle200()
	var status response.StatusCode = response.OK

	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		status = response.BadRequest
		html = handle400()
	case "/myproblem":
		status = response.InternalServerError
		html = handle500()
	}

	html = strings.Trim(html, "\n")

	w.WriteStatusLine(status)
	headers := response.GetDefaultHeaders(len(html))
	headers.Replace("Content-Type", "text/html")
	w.WriteHeaders(headers)

	w.Write([]byte(html))
}

func handle400() string {
	return `
<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>
`
}

func handle500() string {
	return `
<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>
`
}

func handle200() string {
	return `
<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>
`
}
