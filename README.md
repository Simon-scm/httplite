# httplite

this is a small educational HTTP request parser written in Go.

It is intended for my own learning purposes and is by no means production ready.

## Features

- Parses a basic HTTP request line.
- Parses headers into `map[string][]string`.
- Reads request bodies via `Content-Length`.
- Rejects `Transfer-Encoding: chunked`.
- Limits request bodies to 1 MiB.

## Install

```powershell
go get github.com/Simon-scm/httplite
```

## Usage

```go
package main

import (
	"log"
	"net"

	"github.com/Simon-scm/httplite"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go func() {
			defer conn.Close()

			req, err := httplite.ParseRequest(conn)
			if err != nil {
				log.Println(err)
				return
			}

			log.Printf("%s %s %s", req.Method, req.Target, req.Proto)
		}()
	}
}
```

## Example Server

This repository includes a tiny test server:

```powershell
go run ./cmd/httplite-server
```

Then send a request:

```powershell
curl http://localhost:8080/test
```

Or send a POST body:

```powershell
curl -X POST http://localhost:8080/echo -d "hello world"
```

## Tests

```powershell
go test ./...
```

## Limitations

This package currently does not implement the full HTTP specification.

Known limitations:

- No chunked request body support
- Minimal request-line validation
- Minimal header validation
- No keep-alive request loop
- No response writer abstraction
- Not production-ready

