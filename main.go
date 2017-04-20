package main

import (
	"fmt"
	"net"
	"strings"
	"bufio"
)

func game(conn1 net.Conn, conn2 net.Conn, conn1_q chan string, conn2_q chan string) {

	writer1 := bufio.NewWriter(conn1)
	reader1 := bufio.NewReader(conn1)

	writer2 := bufio.NewWriter(conn2)
	reader2 := bufio.NewReader(conn2)

	for {
		// conn1 からの入力待ち
		line1, _ := reader1.ReadString('\n')
		println("conn1 message: ", line1)

		// conn2 からの入力待ち
		line2, _ := reader2.ReadString('\n')
		println("conn2 message: ", line2)

		// conn1 にフレーム情報を返す
		writer1.WriteString(line1 + "#" + line2)
		writer1.Flush()

		// conn2 にフレーム情報を返す
		writer2.WriteString(line1 + "#" + line2)
		writer2.Flush()
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

			if waitConn == nil {
				waitConn = conn
				wait_conn_q = conn_q
			} else {
				go game(waitConn, conn, wait_conn_q, conn_q)
				waitConn = nil
			}

			// 終わるまで待つ
			<-conn_q
		}()
	}
}