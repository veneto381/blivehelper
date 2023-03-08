package blivehelper

import (
	"bytes"
	"encoding/binary"
	"unsafe"
)

type Head struct {
	TotalLen uint32
	HeadLen  uint16
	Version  uint16
	Code     uint32
	Sequence uint32
}

func encode(code int, body []byte) ([]byte, error) {
	head := Head{TotalLen: 0, Version: 1, Code: uint32(code), Sequence: 1}
	head.HeadLen = uint16(unsafe.Sizeof(head))
	head.TotalLen = uint32(head.HeadLen) + uint32(len(body))
	var bufer bytes.Buffer
	err := binary.Write(&bufer, binary.BigEndian, head)
	if err != nil {
		return nil, err
	}
	err = binary.Write(&bufer, binary.BigEndian, body)
	if err != nil {
		return nil, err
	}
	return bufer.Bytes(), nil
}
