package immutable

import (
	"encoding/binary"
	"hash/maphash"
)

var EnableHashCollision = false

var seed = maphash.MakeSeed()

func hash(key any) uint32 {
	switch k := key.(type) {
	case string:
		return uint32(maphash.String(seed, k))
	case int8:
		return uint32(k)
	case uint8:
		return uint32(k)
	case int16:
		return uint32(k)
	case uint16:
		return uint32(k)
	case int32:
		return uint32(k)
	case uint32:
		return k
	case int64:
		return uint32(k)
	case uint64:
		return uint32(k)
	case int:
		return uint32(k)
	case uint:
		return uint32(k)
	case []byte:
		if len(k) == 4 && EnableHashCollision {
			// special case to make colliding hash in order to force a deep tree
			return binary.LittleEndian.Uint32(k)
		}
		return uint32(maphash.Bytes(seed, k))
	default:
		panic(UnhashableKeyType)
	}
}
