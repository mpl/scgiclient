package scgiclient

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
)

// TODO(mpl): I don't like this api. redo it similar to http?

func Receive(conn net.Conn) ([]byte, error) {
	resp := bufio.NewReader(conn)
	terminator := string([]byte{13, 10})
	status, err := resp.ReadString('\n')
	if err != nil {
		return nil, err
	}
	status = strings.TrimRight(status, terminator)
	if status != "Status: 200 OK" {
		return nil, fmt.Errorf("Got %v as response status", status)
	}
	// TODO(mpl): other header fields should not be in the body
	var body bytes.Buffer
	if _, err = io.Copy(&body, resp); err != nil {
		return nil, fmt.Errorf("Could not read response: %v", err)
	}
	return body.Bytes(), nil
}

func Send(addr string, r io.Reader) (net.Conn, error) {
	// TODO(mpl): maybe take the dial out of here and then
	// we do not have to return the conn ?
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	req := newRequest(body)
	if _, err = req.send(conn); err != nil {
		return nil, err
	}
	return conn, nil
}

type request struct {
	header []byte
	body   []byte
}

func newRequest(body []byte) request {
	return request{
		header: defaultHeader(len(body)),
		body:   body,
	}
}

func (r *request) send(c net.Conn) (int64, error) {
	msg := append(netstring(r.header), r.body...)
	return io.Copy(c, bytes.NewReader(msg))
}

// TODO(mpl): report hoisie his scgi server panics if field missing
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
