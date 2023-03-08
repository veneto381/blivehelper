package blivehelper

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"runtime"
	"strings"
)

func MsgDecode(msg []byte) ([][]byte, error) {
	defer func(msg []byte) {
		if err := recover(); err != nil {
			trace := func(message string) string {
				var pcs [32]uintptr
				n := runtime.Callers(3, pcs[:])

				var str strings.Builder
				str.WriteString(message + "\nTraceback:")
				for _, pc := range pcs[:n] {
					fn := runtime.FuncForPC(pc)
					file, line := fn.FileLine(pc)
					str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
				}
				return str.String()
			}
			message := fmt.Sprintf("%s", err)
			log.Printf("%s\n\n", trace(message))
			log.Printf("当前数据：%s", msg)
		}
	}(msg)

	ans := [][]byte{}
	for len(msg) > 15 {
		var totalLen uint32
		var temp bytes.Buffer
		buffer := bytes.NewBuffer(msg[:4])
		err := binary.Read(buffer, binary.BigEndian, &totalLen)
		if err != nil {
			return [][]byte{}, err
		}
		temp.Write(msg[16:totalLen])
		ans = append(ans, temp.Bytes())
		msg = msg[totalLen:]
	}
	return ans, nil
}
