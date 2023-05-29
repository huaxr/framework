// Author: XinRui Hua
// Time:   2022/4/19 上午10:08
// Git:    huaxr

package toolutil

import (
	"bytes"
	"encoding/binary"
)

func Int64ToBytes(data int64) []byte {
	bytebuf := bytes.NewBuffer([]byte{})
	binary.Write(bytebuf, binary.BigEndian, data)
	return bytebuf.Bytes()
}

func Int8ToBytes(data int8) []byte {
	bytebuf := bytes.NewBuffer([]byte{})
	binary.Write(bytebuf, binary.BigEndian, data)
	return bytebuf.Bytes()
}

func BytesToInt64(bys []byte) int64 {
	bytebuff := bytes.NewBuffer(bys)
	var data int64
	binary.Read(bytebuff, binary.BigEndian, &data)
	return data
}

func BytesToInt8(bys []byte) int8 {
	bytebuff := bytes.NewBuffer(bys)
	var data int8
	binary.Read(bytebuff, binary.BigEndian, &data)
	return data
}
