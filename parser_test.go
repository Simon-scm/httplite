package httplite

import (
	"strings"
	"testing"
)

func TestParseRequestPostWithBody(t *testing.T) {
	raw := "POST /echo HTTP/1.1\r\n" +
		"Host: localhost\r\n" +
		"User-Agent: testclient\r\n" +
		"Content-Length: 11\r\n" +
		"\r\n" +
		"hello world"

	req, err := ParseRequest(strings.NewReader(raw))
	if err != nil {
		t.Fatalf("ParseRequest returned error: %v", err)
	}

	if req.Method != "POST" {
		t.Fatalf("Method = %q, want POST", req.Method)
	}
	if req.Target != "/echo" {
		t.Fatalf("Target = %q, want /echo", req.Target)
	}
	if req.Proto != "HTTP/1.1" {
		t.Fatalf("Proto = %q, want HTTP/1.1", req.Proto)
	}
	if got := req.Headers["host"][0]; got != "localhost" {
		t.Fatalf("Host = %q, want localhost", got)
	}
	if got := string(req.Body); got != "hello world" {
		t.Fatalf("Body = %q, want hello world", got)
	}
}

func TestParseRequestMissingHost(t *testing.T) {
	raw := "GET / HTTP/1.1\r\n\r\n"

	_, err := ParseRequest(strings.NewReader(raw))
	if err == nil {
		t.Fatal("ParseRequest returned nil error, want missing host error")
	}
}

func TestParseRequestRejectsChunked(t *testing.T) {
	raw := "POST /upload HTTP/1.1\r\n" +
		"Host: localhost\r\n" +
		"Transfer-Encoding: chunked\r\n" +
		"\r\n"

	_, err := ParseRequest(strings.NewReader(raw))
	if err == nil {
		t.Fatal("ParseRequest returned nil error, want unsupported chunked error")
	}
}

func TestParseRequestInvalidContentLength(t *testing.T) {
	raw := "POST /echo HTTP/1.1\r\n" +
		"Host: localhost\r\n" +
		"Content-Length: nope\r\n" +
		"\r\n"

	_, err := ParseRequest(strings.NewReader(raw))
	if err == nil {
		t.Fatal("ParseRequest returned nil error, want invalid content-length error")
	}
}

func TestParseRequestBodyTooLarge(t *testing.T) {
	raw := "POST /echo HTTP/1.1\r\n" +
		"Host: localhost\r\n" +
		"Content-Length: 1048577\r\n" +
		"\r\n"

	_, err := ParseRequest(strings.NewReader(raw))
	if err == nil {
		t.Fatal("ParseRequest returned nil error, want body too large error")
	}
}
