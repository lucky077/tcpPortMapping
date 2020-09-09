package util

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
)

func ToBytes(head byte, body string) []byte {
	buf := make([]byte, len(body)+5)
	binary.BigEndian.PutUint32(buf, uint32(len(body)+1))
	buf[4] = head

	copy(buf[5:], body)

	return buf
}
func GetData(data []byte) (byte, string) {
	return data[0], string(data[1:])
}

func ErrCheck(err error) {
	if err == nil {
		return
	}
	fmt.Println(err.Error())
}

func GetInput() string {
	in := bufio.NewReader(os.Stdin)
	str, _, err := in.ReadLine()
	if err != nil {
		return err.Error()
	}
	return string(str)
}
