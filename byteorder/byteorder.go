package byteorder

import "encoding/binary"

type ByteOrder struct {
	Uint32 func([]byte) uint32
}

var LittleEndian = ByteOrder{
	Uint32: func(k []byte) uint32 {
		switch len(k) {
		case 0:
			return 0
		case 1:
			return uint32(k[0])
		case 2:
			return uint32(k[0]) | uint32(k[1])<<8
		case 3:
			return uint32(k[0]) | uint32(k[1])<<8 | uint32(k[2])<<16
		default:
			return binary.LittleEndian.Uint32(k)
		}
	},
}

var BigEndian = ByteOrder{
	Uint32: func(k []byte) uint32 {
		switch len(k) {
		case 0:
			return 0
		case 1:
			return uint32(k[0])
		case 2:
			return uint32(k[0])<<8 | uint32(k[1])
		case 3:
			return uint32(k[0])<<16 | uint32(k[1])<<8 | uint32(k[2])
		default:
			return binary.BigEndian.Uint32(k)
		}
	},
}
