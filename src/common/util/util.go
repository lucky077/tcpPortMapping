package util

import (
	"bufio"
	"fmt"
	"os"
)

func ToBytes(head byte, body string) []byte {
	buf := make([]byte, len(body)+1)

	buf[0] = head

	copy(buf[1:], body)

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
