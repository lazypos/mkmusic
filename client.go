package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"sync"
	"time"
)

var errCounts = 0
var totalCounts = 0
var mux sync.Mutex

func main() {
	ip, err := ioutil.ReadFile("ip.txt")
	if err != nil {
		log.Println("读取ip.txt失败", err)
		ip = []byte("127.0.0.1")
	}
	for j := 0; j < 100; j++ {
		log.Println("第", j, "轮")
		wg := &sync.WaitGroup{}
		wg.Add(100)
		for i := 0; i < 100; i++ {
			go ConnectTest(string(ip[:]), wg)
		}
		wg.Wait()
		time.Sleep(time.Second)
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

func ConnectTest(ip string, w *sync.WaitGroup) {
	defer w.Done()
	conn, err := net.Dial("tcp", fmt.Sprint(ip, ":2368"))
	if err != nil {
		log.Println("连接服务器失败", err)
		return
	}
	defer conn.Close()

	err = sendMesssage(conn, []byte("a"))
	if err != nil {
		log.Println("发送消息失败", err)
		return
	}
	_, err = recvMessage(conn, 1)
	if err != nil {
		log.Println("接受消息失败", err)
		return
	}
}
