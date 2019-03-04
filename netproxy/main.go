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

	log.Println("prefix:", prefix)

	remote, err := net.Dial("tcp", "localhost:8081")
	if err != nil {
		log.Println(err)
		return
	}
	defer remote.Close()

	go io.Copy(remote, io.MultiReader(&buf, conn))
	io.Copy(conn, remote)
}
