package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/atotto/clipboard"
	"httpPortMapping/src/common/log"
	"httpPortMapping/src/common/util"
	"net"
	"strconv"
	"strings"
	"sync"
)

const (
	LocalHost = "127.0.0.1"
	//LocalHost  = "192.168.1.106"
	//RemoteHost  = LocalHost
	RemoteHost = "118.25.13.237"
	Port       = 28999
	HeadSize   = 4
	BufSize    = 2048
	BufSizeUdp = 65535
)

var (
	port string

	localPort string
	network   string

	rlMap = sync.Map{}
	lrMap = sync.Map{}

	pcMap = sync.Map{}
)

func main() {
	log.Init()
	run()
}

func run() {
	fmt.Println("enter tcp/udp + port")
	fmt.Println("demo : tcp 8080")

	row := util.GetInput()

	fields := strings.Fields(row)
	network = fields[0]
	localPort = fields[1]

	connectTCPCoder(address(RemoteHost, Port), func(conn net.Conn) {
		conn.Write(util.ToBytes(1, "hello"))
		conn.Write(util.ToBytes(2, network))
	}, func(conn net.Conn, data []byte) {

		head, body := util.GetData(data)

		//fmt.Println(head,body)

		switch head {
		case 1:

			break
		case 2:
			port = body
			tip := RemoteHost + ":" + port
			fmt.Println(tip + "  已经复制到剪贴板")
			clipboard.WriteAll(tip)
			break
		case 21:
			conn, _ := pcMap.Load(body)
			if conn != nil {
				conn.(net.Conn).Close()
				conn0, _ := lrMap.Load(conn)
				if conn0 != nil {
					conn0.(net.Conn).Close()
				}
			}
			break
		case 11:
			remote := make(chan net.Conn)
			local := make(chan net.Conn)
			go connectTCP(addressStr(LocalHost, localPort), func(conn net.Conn) {
				local <- conn

				pcMap.Store(body, conn)

				lrMap.Store(conn, <-remote)

			}, func(conn net.Conn, data []byte) {
				//fmt.Println(string(data))

				conn0, _ := lrMap.Load(conn)
				conn0.(net.Conn).Write(data)
			})

			go connectTCP(addressStr(RemoteHost, body), func(conn net.Conn) {
				rlMap.Store(conn, <-local)

				remote <- conn
			}, func(conn net.Conn, data []byte) {
				//fmt.Println(string(data))

				conn0, _ := rlMap.Load(conn)
				conn0.(net.Conn).Write(data)
			})

			break
		case 12:
			remote := make(chan net.Conn)
			local := make(chan net.Conn)
			go connectUDP(addressStr(LocalHost, localPort), func(conn net.Conn) {
				local <- conn

				pcMap.Store(body, conn)

				lrMap.Store(conn, <-remote)

			}, func(conn net.Conn, data []byte) {
				//fmt.Println(string(data))

				conn0, _ := lrMap.Load(conn)
				conn0.(net.Conn).Write(data)
			})

			go connectUDP(addressStr(RemoteHost, body), func(conn net.Conn) {
				rlMap.Store(conn, <-local)

				remote <- conn

				conn.Write([]byte(port))
			}, func(conn net.Conn, data []byte) {
				//fmt.Println(string(data))

				conn0, _ := rlMap.Load(conn)
				conn0.(net.Conn).Write(data)
			})
			break
		}

	})

	util.GetInput()
}

func connectTCPCoder(addr string, cb func(conn net.Conn), cb2 func(conn net.Conn, data []byte)) {

	conn, e := net.Dial("tcp", addr)
	util.ErrCheck(e)
	cb(conn)

	var (
		buffer    = bytes.NewBuffer(make([]byte, 0, BufSize)) //buffer用来缓存读取到的数据
		readBytes = make([]byte, BufSize)                     //readBytes用来接收每次读取的数据，每次读取完成之后将readBytes添加到buffer中
		isHead    = true                                      //用来标识当前的状态：正在处理size部分还是body部分
		bodyLen   = 0                                         //表示body的长度
	)

	for {
		//首先读取数据
		readByteNum, err := conn.Read(readBytes)
		if err != nil {
			util.ErrCheck(err)
			break
		}
		buffer.Write(readBytes[0:readByteNum]) //将读取到的数据放到buffer中

		// 然后处理数据
		for {
			if isHead {
				if buffer.Len() >= HeadSize {
					isHead = false
					head := make([]byte, HeadSize)
					_, err = buffer.Read(head)
					util.ErrCheck(err)
					bodyLen = int(binary.BigEndian.Uint32(head))
				} else {
					break
				}
			}

			if !isHead {
				if buffer.Len() >= bodyLen {
					body := make([]byte, bodyLen)
					_, err = buffer.Read(body[:bodyLen])
					util.ErrCheck(err)
					cb2(conn, body[:bodyLen])
					isHead = true
				} else {
					break
				}
			}
		}
	}

}

func connectTCP(addr string, cb func(conn net.Conn), cb2 func(conn net.Conn, data []byte)) {

	conn, err := net.Dial("tcp", addr)

	util.ErrCheck(err)

	cb(conn)

	defer func() {
		if recover() != nil {
			conn0, _ := rlMap.Load(conn)
			if conn0 != nil {
				conn0.(net.Conn).Close()
			} else {
				conn0, _ := lrMap.Load(conn)
				if conn0 != nil {
					conn0.(net.Conn).Close()
				}
			}
		}
	}()

	readBytes := make([]byte, BufSize)

	for {
		n, err := conn.Read(readBytes)
		if err != nil {
			util.ErrCheck(err)

			if err.Error() == "EOF" {
				connNew, err := net.Dial("tcp", addr)
				util.ErrCheck(err)

				conn0, _ := lrMap.Load(conn)
				conn = connNew

				if conn0 != nil {
					lrMap.Store(connNew, conn0)
					rlMap.Store(conn0, connNew)
				}

				continue
			} else {
				break
			}

		}

		cb2(conn, readBytes[:n])
	}

}

func connectUDP(addr string, cb func(conn net.Conn), cb2 func(conn net.Conn, data []byte)) {

	conn, err := net.Dial("udp", addr)
	util.ErrCheck(err)
	cb(conn)

	for {
		buf := make([]byte, BufSizeUdp)

		n, err := conn.Read(buf)
		util.ErrCheck(err)

		cb2(conn, buf[:n])
	}

}

func address(host string, port int) string {
	return host + ":" + strconv.Itoa(port)
}

func addressStr(host string, port string) string {
	return host + ":" + port
}
