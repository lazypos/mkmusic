package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

func main() {
	listenfd, err := net.Listen("tcp", ":2368")
	if err != nil {
		log.Println("侦听失败.", err)
		return
	}
	defer listenfd.Close()

	for {
		if conn, err := listenfd.Accept(); err != nil {
			log.Println(`accept错误`, err)
			time.Sleep(time.Second)
		} else {
			log.Println(`收到连接`, conn.RemoteAddr().String())
			conn.(*net.TCPConn).SetLinger(0)
			go processConn(conn)
		}
	}
}

//基本发送函数
func sendMesssage(conn net.Conn, msg []byte) error {
	sendlen := 0
	for sendlen < len(msg) {
		n, err := conn.Write(msg[sendlen:])
		if err != nil {
			return fmt.Errorf(`[NET] sendMesssage error: %v`, err)
		} else if n <= 0 {
			return fmt.Errorf(`[NET] sendMesssage error: remote closed`)
		}
		sendlen += n
	}
	return nil
}

//基本接收函数
func recvMessage(conn net.Conn, totallen int) ([]byte, error) {
	recvlen := 0
	buf := make([]byte, totallen)
	for recvlen < totallen {
		n, err := conn.Read(buf[recvlen:])
		if err != nil {
			return []byte{}, fmt.Errorf(`[NET] recvMessage error: %v`, err)
		} else if n <= 0 {
			return []byte{}, fmt.Errorf(`[NET] recvMessage error: remote closed`)
		}
		recvlen += n
	}
	return buf, nil
}

var errCounts = 0
var totalCounts = 0
var mux sync.Mutex

func processConn(conn net.Conn) {
	now := time.Now().Unix()
	totalCounts += 1
	defer conn.Close()
	c, err := recvMessage(conn, 1)
	if err != nil {
		log.Println("接收错误", err, conn.RemoteAddr())
		errCounts += 1
	} else {
		e := sendMesssage(conn, c)
		if e != nil {
			log.Println("发送错误", conn.RemoteAddr())
			errCounts += 1
		}
	}
	recvMessage(conn, 1)
	log.Println("处理完毕，耗时", time.Now().Unix()-now, "秒, 总失败:", errCounts)
}
