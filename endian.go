package gtp5gnl

import (
	"encoding/binary"
	"unsafe"
)

var native binary.ByteOrder = NativeEndian()

func NativeEndian() binary.ByteOrder {
	var x uint32 = 0x01020304
	if *(*byte)(unsafe.Pointer(&x)) == 0x01 {
		return binary.BigEndian
	} else {
		return binary.LittleEndian
	}
}
