/*
Copyright 2013 Mathieu Lonjaret.
*/

package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mpl/scgiclient"
)

func ghettoXMLRpc(command, arg string) string {
	// TODO(mpl): allow more args, types other than string maybe.
	xml := `<?xml version="1.0"?>
	<methodCall>
		<methodName>` + command + `</methodName>`
	if arg != "" {
		xml += `<params><param><value><string>` + arg + `</string></value></param></params>`
	}
	xml += "</methodCall>"
	return xml
}

func usage() {
	fmt.Print("	usage: rtorrentrpc host:port rpccommand [arg]\n")
	fmt.Print("	See http://libtorrent.rakshasa.no/wiki/RTorrentCommands for the list of commands.\n")
	os.Exit(1)
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		usage()
	}
	addr := args[0]
	command := args[1]
	cmdArg := ""
	if args[2] != "" {
		cmdArg = args[2]
	}
	xmlrpc := ghettoXMLRpc(command, cmdArg)
	resp, err := scgiclient.Send(addr, bytes.NewReader([]byte(xmlrpc)))
	if err != nil {
		log.Fatal(err)
	}
	// TODO(mpl): Unmarshall response. meh.
	fmt.Printf("%v", string(resp.Body))
}
