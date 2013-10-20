package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"

	"github.com/mpl/scgiclient"
)

func ughxml(command string) string {
	return `<?xml version="1.0"?>
	<methodCall>
		<methodName>` + command + `</methodName>
	</methodCall>`
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		log.Fatal("wat")
	}
	// TODO(mpl): hostport parsing. trust for now.
	addr := args[0]
	xmlrpc := ughxml(args[1])
	conn, err := scgiclient.Send(addr, bytes.NewReader([]byte(xmlrpc)))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	resp, err := scgiclient.Receive(conn)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v", string(resp))
}
