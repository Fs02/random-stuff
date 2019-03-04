package main

import (
	"bytes"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("connected:", conn.RemoteAddr())
		go proxy(conn)
	}
}

func print(conn net.Conn) {
	defer conn.Close()
	io.Copy(os.Stdout, conn)
}

var remotes = map[string]string{
	"hello": "localhost:8081",
	"other": "localhost:8082",
}

func proxy(conn net.Conn) {
	defer conn.Close()

	var buf bytes.Buffer
	io.CopyN(&buf, conn, 15)
	log.Println("peek:", buf.String())

	prefix := ""
	parts := strings.Split(buf.String(), " ")
	if len(parts) > 1 {
		paths := strings.Split(parts[1], "/")
		if len(paths) > 1 {
			prefix = paths[1]
		}
	}

	if remotes[prefix] == "" {
		conn.Write([]byte("HTTP/1.1 404 Not found\r\nContent-Length: 0\r\n\r\n"))
		return
	}

	remote, err := net.Dial("tcp", remotes[prefix])
	if err != nil {
		log.Println(err)
		return
	}
	defer remote.Close()

	go io.Copy(remote, io.MultiReader(&buf, conn))
	io.Copy(conn, remote)
}
