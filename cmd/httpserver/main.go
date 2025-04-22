package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"httpfromtcp/internal/headers"
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
	h := response.GetDefaultHeaders(len(html))

	if strings.HasPrefix(target, "/httpbin") {
		h.Remove("Content-Length")
		h.Set("Transfer-Encoding", "chunked")
		h.Set("Trailer", "X-Content-Sha256")
		h.Set("Trailer", "X-Content-Length")
		w.WriteHeaders(h)

		path := strings.TrimPrefix(target, "/httpbin/")
		url := fmt.Sprintf("https://httpbin.org/%s", path)

		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}

		var bodyBuf bytes.Buffer
		buf := make([]byte, 1024)

		for {
			n, _ := resp.Body.Read(buf)
			if n == 0 {
				break
			}
			w.WriteChunkedBody(buf[:n])
			bodyBuf.Write(buf[:n])
		}

		w.WriteChunkedBodyDone()

		fullBody := bodyBuf.Bytes()
		trailers := headers.NewHeaders()
		sha256 := fmt.Sprintf("%x", sha256.Sum256(fullBody))
		trailers.Set("X-Content-SHA256", sha256)
		trailers.Set("X-Content-Length", fmt.Sprintf("%d", len(fullBody)))
		err = w.WriteTrailers(trailers)
		if err != nil {
			fmt.Println("Error writing trailers:", err)
		}
		fmt.Println("Wrote trailers")

		return
	}

	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
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
