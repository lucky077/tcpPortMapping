package main

import (
	"fmt"
	"httpPortMapping/src/common/log"
	"httpPortMapping/src/common/util"
	"net"
	"strconv"
	"sync"
)

type session struct {
	//基本
	client net.TCPConn
	//映射的端口
	port int
}

var (
	clientServer sync.Map
	serverClinet sync.Map
)

func main() {

	log.Init()

	fmt.Println("server run")
	serveTcp(28999, func(conn net.Conn) {
		log.Info(conn.RemoteAddr().String() + " connected")
	}, func(conn net.Conn, data []byte) {
		head, body := util.GetData(data)

		println(head)
		println(body)

	})
}

func serveTcp(port int, cb func(conn net.Conn), cb2 func(conn net.Conn, data []byte)) {

	listener, e := net.Listen("tcp", ":"+strconv.Itoa(port))
	util.ErrCheck(e)

	for {
		conn, e := listener.Accept()
		util.ErrCheck(e)
		if cb != nil {
			cb(conn)
		}

		buf := make([]byte, 1024)

		for {
			n, e := conn.Read(buf)
			if e != nil {
				util.ErrCheck(e)
				break
			}

			go cb2(conn, buf[:n])
		}

	}
}

func serveUdp(port int, cb func(conn net.Conn), cb2 func(conn net.Conn, data []byte)) {

	addr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(port))
	util.ErrCheck(err)

	conn, err := net.ListenUDP("udp", addr)

	util.ErrCheck(err)

	if cb != nil {
		cb(conn)
	}

	buf := make([]byte, 1024)

	for {
		n, e := conn.Read(buf)
		if e != nil {
			util.ErrCheck(e)
			break
		}

		go cb2(conn, buf[:n])
	}

}
