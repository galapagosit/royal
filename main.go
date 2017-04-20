package main

import (
	"fmt"
	"net"
	"strings"
)

func game(conn1 net.Conn, conn2 net.Conn, conn1_q chan string, conn2_q chan string) {
	for {
		// conn1 からの入力待ち
		message1 := <-conn1_q
		println("conn1 message: ", message1)

		// conn2 からの入力待ち
		message2 := <-conn2_q
		println("conn2 message: ", message2)

		// conn1 にフレーム情報を返す
		conn1.Write([]byte(message1 + "#" + message2 + "\r\n"))

		// conn2 にフレーム情報を返す
		conn2.Write([]byte(message1 + "#" + message2 + "\r\n"))
	}

	// 親のgoroutineを終わらせる
	conn1_q <- "end"
	conn2_q <- "end"
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	var waitConn net.Conn
	var wait_conn_q chan string

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic("accept error")
		}

		go func() {
			fmt.Printf("Accept %v\n", conn.RemoteAddr())

			conn_q := make(chan string)
			var buf string = ""

			if waitConn == nil {
				waitConn = conn
				wait_conn_q = conn_q
			} else {
				go game(waitConn, conn, wait_conn_q, conn_q)
				waitConn = nil
			}

			for {
				messageBuf := make([]byte, 1024)
				messageLen, err := conn.Read(messageBuf)
				if err != nil {
					panic("read error")
				}
				message := string(messageBuf[:messageLen])
				println("message: ", message)

				splited := strings.Split(message, "\r\n")

				buf += splited[len(splited) - 1]
				lines := splited[:len(splited) - 1]

				for _, line := range lines {
					conn_q <- line
				}
			}
		}()
	}
}