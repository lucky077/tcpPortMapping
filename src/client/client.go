package main

import (
	"fmt"
	"httpPortMapping/src/common/util"
	"net"
)

const (
	ADDRESS = "127.0.0.1" + ":28999"
)

func main() {
	fmt.Println("client run")

	connect(ADDRESS, func(conn net.Conn) {
		conn.Write(util.ToBytes(1, ""))
	})

	util.GetInput()
}

func connect(addr string, cb func(conn net.Conn)) {

	conn, e := net.Dial("tcp", addr)
	util.ErrCheck(e)
	cb(conn)

}
