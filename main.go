package main

import (
	"fmt"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic("accept error")
		}
		go func() {
			fmt.Printf("Accept %v\n", conn.RemoteAddr())
			for {
				messageBuf := make([]byte, 1024)
				messageLen, err := conn.Read(messageBuf)
				if err != nil {
					panic("read error")
				}
				message := string(messageBuf[:messageLen])
				print("message: ", message)
			}
		}()
	}
}