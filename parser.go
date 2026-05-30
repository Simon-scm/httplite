package httplite

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Request struct {
	Method, Target, Proto string
	Headers               map[string][]string
	Body                  []byte
}

func readCRLF(br *bufio.Reader) (string, error) {
	s, err := br.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}

	s = strings.TrimSuffix(s, "\n")
	s = strings.TrimSuffix(s, "\r")

	if err == io.EOF && len(s) == 0 {
		return "", io.EOF
	}
	return s, nil
}

func parseHeaders(br *bufio.Reader) (map[string][]string, error) {
	headers := make(map[string][]string)

	for {
		line, err := readCRLF(br)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if line == "" {
			break
		}

		idx := strings.Index(line, ":")
		if idx < 0 {
			return nil, fmt.Errorf("malformed Header")
		}

		key := strings.ToLower(strings.TrimSpace(line[:idx]))
		val := strings.TrimSpace(line[idx+1:])

		headers[key] = append(headers[key], val)

	}

	if _, ok := headers["host"]; !ok {
		return nil, fmt.Errorf("no host in header")
	}

	return headers, nil
}

func parseRequestLine(br *bufio.Reader) (string, string, string, error) {
	line, err := readCRLF(br)
	if err != nil {
		return "", "", "", err
	}

	parts := strings.SplitN(line, " ", 3)

	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("request line not complete")
	}

	method, target, proto := parts[0], parts[1], parts[2]

	return method, target, proto, nil
}

func readNBytes(br *bufio.Reader, n int) ([]byte, error) {
	buf := make([]byte, n)
	_, err := io.ReadFull(br, buf)
	return buf, err
}

// maxBody: Sicherheits-Limit in Bytes
func parseBodyCL(br *bufio.Reader, headers map[string][]string, maxBody int) ([]byte, error) {
	vals, ok := headers["content-length"]
	if !ok || len(vals) == 0 {
		return nil, nil
	}

	clStr := strings.TrimSpace(vals[0])
	n, err := strconv.Atoi(clStr)
	if err != nil || n < 0 {
		return nil, fmt.Errorf("invalid Content-Length: %q", clStr)
	}

	if maxBody > 0 && n > maxBody {
		return nil, fmt.Errorf("body too large: %d > %d", n, maxBody)
	}

	body, err := readNBytes(br, n)
	if err != nil {
		return nil, fmt.Errorf("reading body failed %w", err)
	}

	return body, nil
}

func ParseRequest(r io.Reader) (*Request, error) {
	br, ok := r.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(r)
	}

	method, target, proto, err := parseRequestLine(br)
	if err != nil {
		return nil, err
	}

	headers, err := parseHeaders(br)
	if err != nil {
		return nil, err
	}

	// Transfer encoding chunked aktuell ablehnen (bis implementieren)
	if teVals, ok := headers["transfer-encoding"]; ok && len(teVals) > 0 {
		te := strings.ToLower(strings.TrimSpace(teVals[0]))
		if strings.Contains(te, "chunked") {
			return nil, fmt.Errorf("unsupported transfer-encoding: chunked")
		}
	}

	req := Request{
		Method:  method,
		Target:  target,
		Proto:   proto,
		Headers: headers,
	}

	const maxBody = 1 << 20
	body, err := parseBodyCL(br, headers, maxBody)
	if err != nil {
		return nil, err
	}
	req.Body = body

	return &req, nil
}
