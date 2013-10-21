package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"

	"github.com/mpl/scgiclient"
)

func ughxml(command, arg string) string {
	xml := `<?xml version="1.0"?>
	<methodCall>
		<methodName>` + command + `</methodName>`
	if arg != "" {
		xml += `<params><param><value><string>` + arg + `</string></value></param></params>`
	}
	xml += "</methodCall>"
	return xml
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		log.Fatal("wat")
	}
	// TODO(mpl): hostport parsing. trust for now.
	addr := args[0]
	command := args[1]
	cmdArg := ""
	if args[2] != "" {
		cmdArg = args[2]
	}
	xmlrpc := ughxml(command, cmdArg)
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
