package main

import (
	"fmt"
	"net"
	"strings"
)

func game(conn1 net.Conn, conn2 net.Conn, conn1_q chan struct{}, conn2_q chan struct{}) {

	// 各クライアントの0フレーム目開始の合図
	conn1.Write([]byte("start"))
	conn2.Write([]byte("start"))

	for {
		// conn1 からの入力待ち
		messageBuf1 := make([]byte, 1024)
		messageLen1, err1 := conn1.Read(messageBuf1)
		if err1 != nil {
			panic("conn1 read error")
		}
		message1 := string(messageBuf1[:messageLen1])
		println("conn1 message: ", message1)

		// conn2 からの入力待ち
		messageBuf2 := make([]byte, 1024)
		messageLen2, err2 := conn2.Read(messageBuf2)
		if err2 != nil {
			panic("conn2 read error")
		}
		message2 := string(messageBuf2[:messageLen2])
		println("conn2 message: ", message2)

		// conn1 にフレーム情報を返す
		conn1.Write([]byte(message1 + "#" + message2))

		// conn2 にフレーム情報を返す
		conn2.Write([]byte(message1 + "#" + message2))
	}

	// 親のgoroutineを終わらせる
	conn1_q <- struct{}{}
	conn2_q <- struct{}{}
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	var waitConn net.Conn
	var wait_conn_q chan struct{}

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic("accept error")
		}

		go func() {
			fmt.Printf("Accept %v\n", conn.RemoteAddr())

			conn_q := make(chan struct{})

			messageBuf := make([]byte, 256)
			messageLen, err := conn.Read(messageBuf)
			if err != nil {
				panic("read error")
			}
			message := string(messageBuf[:messageLen])
			println("message: ", message)
			if strings.TrimRight(message, "\r\n") == "ready" {
				if waitConn == nil {
					waitConn = conn
					wait_conn_q = conn_q
				} else {
					go game(waitConn, conn, wait_conn_q, conn_q)
					waitConn = nil
				}
			}

			// gameが終了するまで待つ
			<-conn_q
		}()
	}
}