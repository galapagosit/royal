package main

import (
	"fmt"
	"net"
	"bufio"
	"strings"
)

func game(conn1 net.Conn, conn2 net.Conn) {

	writer1 := bufio.NewWriter(conn1)
	reader1 := bufio.NewReader(conn1)

	writer2 := bufio.NewWriter(conn2)
	reader2 := bufio.NewReader(conn2)

	for {
		// conn1 からの入力待ち
		line1, err1 := reader1.ReadString('\n')
		if err1 != nil {
			println("conn1 read error")
			return
		}
		line1 = strings.TrimRight(line1, "\r\n")
		println("conn1 message:", line1)

		// conn2 からの入力待ち
		line2, err2 := reader2.ReadString('\n')
		if err2 != nil {
			println("conn2 read error")
			return
		}
		line2 = strings.TrimRight(line2, "\r\n")
		println("conn2 message:", line2)

		// conn1 にフレーム情報を返す
		var message string = line1 + "#" + line2 + "\r\n"

		writer1.WriteString(message)
		writer1.Flush()
		println("conn1 send:" + message)

		// conn2 にフレーム情報を返す
		writer2.WriteString(message)
		writer2.Flush()
		println("conn2 send:" + message)
	}
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	for {
		conn1, err1 := ln.Accept()
		if err1 != nil {
			panic("conn1 accept error")
		}
		fmt.Printf("Accept %v\n", conn1.RemoteAddr())

		conn2, err2 := ln.Accept()
		if err2 != nil {
			panic("conn2 accept error")
		}
		fmt.Printf("Accept %v\n", conn2.RemoteAddr())

		fmt.Printf(">>> game start")
		go game(conn1, conn2)
	}
}
