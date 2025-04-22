package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"log"
	"net/http"
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

	target := req.RequestLine.RequestTarget
	switch target {
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

	if strings.HasPrefix(target, "/httpbin") {
		headers.Delete("Content-Length")
		headers.Set("Transfer-Encoding", "chunked")
		w.WriteHeaders(headers)

		path := strings.TrimPrefix(target, "/httpbin/")
		url := fmt.Sprintf("https://httpbin.org/%s", path)

		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}

		buf := make([]byte, 1024)

		for {
			n, _ := resp.Body.Read(buf)
			if n == 0 {
				break
			}
			w.WriteChunkedBody(buf[0:n])
		}

		w.WriteChunkedBodyDone()
		return
	}

	headers.Replace("Content-Type", "text/html")
	w.WriteHeaders(headers)
	w.WriteBody([]byte(html))
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
