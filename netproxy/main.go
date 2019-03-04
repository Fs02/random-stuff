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
		go request(conn)
	}
}

func print(conn net.Conn) {
	defer conn.Close()
	io.Copy(os.Stdout, conn)
}

func request(conn net.Conn) {
	defer conn.Close()

	remote, err := net.Dial("tcp", "localhost:8081")
	if err != nil {
		log.Println(err)
		return
	}
	defer remote.Close()

	remote.Write([]byte("GET / HTTP/1.0\r\n\r\n"))
	io.Copy(conn, remote)
}
