package blivehelper

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"io"
)

func Decode(data []byte) ([]byte, error) {
	var totalLen uint32
	var headLen uint16
	var version uint16
	var err error

	buffer := bytes.NewReader(data[:4])
	err = binary.Read(buffer, binary.BigEndian, &totalLen)
	if err != nil {
		return []byte{}, err
	}
	buffer = bytes.NewReader(data[4:6])
	err = binary.Read(buffer, binary.BigEndian, &headLen)
	if err != nil {
		return []byte{}, err
	}
	buffer = bytes.NewReader(data[6:8])
	err = binary.Read(buffer, binary.BigEndian, &version)
	if err != nil {
		return []byte{}, err
	}
	buffer = bytes.NewReader(data[headLen:totalLen])

	if version == 2 {
		reader, err := zlib.NewReader(buffer)
		if err != nil {
			return []byte{}, err
		}
		debuff := new(bytes.Buffer)
		_, err = io.Copy(debuff, reader)
		reader.Close()
		if err != nil {
			return []byte{}, err
		}
		return debuff.Bytes(), nil
	} else {
		msg, err := io.ReadAll(buffer)
		if err != nil {
			return []byte{}, err
		}
		return msg, nil
	}
}
