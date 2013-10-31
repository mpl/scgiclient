/*
Copyright 2013 Mathieu Lonjaret.
*/

// Package scgiclient implements the client side of the
// Simple Common Gateway Interface protocol, as described
// at http://python.ca/scgi/protocol.txt
package scgiclient

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
)

// Send sends and scgi request to addr, with
// all the data read from r as the body of the request.
// The received response is returned and the connection
// to addr is closed.
func Send(addr string, r io.Reader) (*Response, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	req := NewRequest(addr, body)
	req.conn = conn
	if _, err = req.Send(); err != nil {
		return nil, err
	}
	return req.Receive()
}

type Request struct {
	Addr   string
	Header []byte
	Body   []byte
	conn   net.Conn
}

// NewRequest returns a Request ready to be sent with
// Request.Send().
func NewRequest(addr string, body []byte) *Request {
	return &Request{
		Addr:   addr,
		Header: defaultHeader(len(body)),
		Body:   body,
	}
}

// Close closes the connection to r.Addr.
func (r *Request) Close() error {
	return r.conn.Close()
}

// Send sends an scgi message built with
// r.Header and r.Body, to r.Addr. It returns the
// number of bytes sent and an error, if any.
func (r *Request) Send() (int64, error) {
	var err error
	if r.conn == nil {
		r.conn, err = net.Dial("tcp", r.Addr)
		if err != nil {
			return 0, err
		}
	}
	msg := append(netstring(r.Header), r.Body...)
	return io.Copy(r.conn, bytes.NewReader(msg))
}

// Receive reads the response from r.Addr and returns
// it in a Response. The connection to r.Addr must already
// be established.
func (r *Request) Receive() (*Response, error) {
	if r.conn == nil {
		return nil, errors.New("Can not receive on a closed connection")
	}
	return receive(r.conn)
}

type Response struct {
	Header map[string]string
	Body   []byte
	conn   net.Conn
}

// Close closes the connection to r.Addr.
func (r *Response) Close() error {
	return r.conn.Close()
}

func receive(conn net.Conn) (*Response, error) {
	r := bufio.NewReader(conn)
	header := make(map[string]string)
	terminator := string([]byte{13, 10})
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return nil, err
		}
		line = strings.TrimRight(line, terminator)
		if line == "" {
			break
		}
		keyValue := strings.SplitN(line, ": ", 2)
		if len(keyValue) != 2 {
			return nil, fmt.Errorf("Bogus header line in response: %q", line)
		}
		header[keyValue[0]] = keyValue[1]
	}

	resp := &Response{
		Header: header,
		conn:   conn,
	}
	status, ok := header["Status"]
	if !ok {
		return nil, errors.New("Did not get a status line in response header")
	}
	if status != "200 OK" {
		return resp, fmt.Errorf("Got %v as response status", status)
	}
	var body bytes.Buffer
	if _, err := io.Copy(&body, r); err != nil {
		return nil, fmt.Errorf("Could not read response: %v", err)
	}
	resp.Body = body.Bytes()
	return resp, nil
}

func defaultHeader(bodyLen int) []byte {
	var dh []byte
	defaultHeaderFields["CONTENT_LENGTH"] = strconv.Itoa(bodyLen)
	for k, v := range defaultHeaderFields {
		dh = append(dh, header(k, v)...)
	}
	return dh
}

func header(name, value string) []byte {
	h := append([]byte(name), 0)
	h = append(h, []byte(value)...)
	return append(h, 0)
}

var defaultHeaderFields = map[string]string{
	"CONTENT_LENGTH":  "",
	"SCGI":            "1",
	"REQUEST_METHOD":  "POST",
	"SERVER_PROTOCOL": "HTTP/1.1",
}

const (
	comma = byte(',')
	colon = byte(':')
)

func netstring(s []byte) []byte {
	le := []byte(strconv.Itoa(len(s)))
	ns := append(le, colon)
	ns = append(ns, s...)
	ns = append(ns, comma)
	return ns
}
