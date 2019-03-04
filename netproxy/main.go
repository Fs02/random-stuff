package main

import (
	"io"
	"log"
	"net"
	"os"
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
		go print(conn)
	}
}

func print(conn net.Conn) {
	defer conn.Close()
	io.Copy(os.Stdout, conn)
}
