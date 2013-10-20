package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	//	"net/http"
	"strconv"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:5000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	err = writeRequest(conn)
	if err != nil {
		log.Fatal(err)
	}
	resp := bufio.NewReader(conn)
	terminator := string([]byte{13, 10})
	status, err := resp.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	status = strings.TrimRight(status, terminator)
	if status != "Status: 200 OK" {
		for _, v := range status {
			println(v)
		}
		log.Fatal(status)
	}
	body := bytes.NewBuffer([]byte{})
	if _, err = io.Copy(body, resp); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v", body.String())
}

var (
	comma = []byte(",")
	colon = []byte(":")
)

func netstring(s []byte) []byte {
	le := []byte(strconv.Itoa(len(s)))
	ns := append(le, byte(':'))
	ns = append(ns, s...)
	ns = append(ns, byte(','))
	return ns
}

func message(header, body []byte) []byte {
	return append(netstring(header), body...)
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

// TODO(mpl): report hoisie it panics if field missing
func defaultHeader(bodyLen int) []byte {
	var dh []byte
	defaultHeaderFields["CONTENT_LENGTH"] = strconv.Itoa(bodyLen)
	for k, v := range defaultHeaderFields {
		dh = append(dh, header(k, v)...)
	}
	return dh
}

var command = `<?xml version="1.0"?>
<methodCall>
	<methodName>get_upload_rate</methodName>
</methodCall>
`

func writeRequest(fd io.ReadWriteCloser) error {
	//	command := "get_upload_rate"
	header := defaultHeader(len(command))
	msg := message(header, []byte(command))
	_, err := io.Copy(fd, bytes.NewReader(msg))
	return err
}
