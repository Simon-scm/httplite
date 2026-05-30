package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Simon-scm/httplite"
)

func main() {
	const addr = ":8080"

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	log.Printf("httplite server listening on http://localhost%s", addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept failed: %v", err)
			continue
		}

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	req, err := httplite.ParseRequest(conn)
	if err != nil {
		log.Printf("parse failed: %v", err)
		fmt.Fprint(conn, "HTTP/1.1 400 Bad Request\r\nConnection: close\r\nContent-Length: 11\r\n\r\nbad request")
		return
	}

	log.Printf("%s %s %s", req.Method, req.Target, req.Proto)

	body := fmt.Sprintf("parsed %s %s\n", req.Method, req.Target)
	fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\nConnection: close\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(body), body)
}
