package misc

func Uint16_Big(b []byte) uint16 {
	return uint16(b[1]) | uint16(b[0])<<8
}

func Uint16_Little(b []byte) uint16 {
	return uint16(b[0]) | uint16(b[1])<<8
}

func Uint24_Big(b []byte) uint32 {
	return uint32(b[2]) | uint32(b[1])<<8 | uint32(b[0])<<16
}
func Uint24_Little(b []byte) uint32 {
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16
}

func Uint32_Big(b []byte) uint32 {
	return uint32(b[3]) | uint32(b[2])<<8 | uint32(b[1])<<16 | uint32(b[0])<<24
}
func Uint32_Little(b []byte) uint32 {
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}

func Uint64_Big(b []byte) uint64 {
	return uint64(b[7]) | uint64(b[6])<<8 | uint64(b[5])<<16 | uint64(b[4])<<24 |
		uint64(b[3])<<32 | uint64(b[2])<<40 | uint64(b[1])<<48 | uint64(b[0])<<56
}

func Uint64_Little(b []byte) uint64 {
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
}

func PutUint16_Little(b []byte, v uint16) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
}

func PutUint16_Big(b []byte, v uint16) {
	b[0] = byte(v >> 8)
	b[1] = byte(v)
}

func PutUint24_Little(b []byte, v uint32) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
}

func PutUint24_Big(b []byte, v uint32) {
	b[0] = byte(v >> 16)
	b[1] = byte(v >> 8)
	b[2] = byte(v)
}

func PutUint32_Little(b []byte, v uint32) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
}

func PutUint32_Big(b []byte, v uint32) {
	b[0] = byte(v >> 24)
	b[1] = byte(v >> 16)
	b[2] = byte(v >> 8)
	b[3] = byte(v)
}

func PutUint64_Little(b []byte, v uint64) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	b[4] = byte(v >> 32)
	b[5] = byte(v >> 40)
	b[6] = byte(v >> 48)
	b[7] = byte(v >> 56)
}

func PutUint64_Big(b []byte, v uint64) {
	b[0] = byte(v >> 56)
	b[1] = byte(v >> 48)
	b[2] = byte(v >> 40)
	b[3] = byte(v >> 32)
	b[4] = byte(v >> 24)
	b[5] = byte(v >> 16)
	b[6] = byte(v >> 8)
	b[7] = byte(v)
}

func Int64ToBytes(n int64) [8]byte {
	b := [8]byte{}
	b[7] = byte(n & 0xff)
	b[6] = byte(n >> 8 & 0xff)
	b[5] = byte(n >> 16 & 0xff)
	b[4] = byte(n >> 24 & 0xff)
	b[3] = byte(n >> 32 & 0xff)
	b[2] = byte(n >> 40 & 0xff)
	b[1] = byte(n >> 48 & 0xff)
	b[0] = byte(n >> 56 & 0xff)
	return b
}

func Int64ToBytesA(n int64, array []byte, offset int) {
	array[7+offset] = byte(n & 0xff)
	array[6+offset] = byte(n >> 8 & 0xff)
	array[5+offset] = byte(n >> 16 & 0xff)
	array[4+offset] = byte(n >> 24 & 0xff)
	array[3+offset] = byte(n >> 32 & 0xff)
	array[2+offset] = byte(n >> 40 & 0xff)
	array[1+offset] = byte(n >> 48 & 0xff)
	array[0+offset] = byte(n >> 56 & 0xff)
}

func BytesToInt64(array []byte) int64 {
	return ((int64(array[0]) & 0xff) << 56) + ((int64(array[1]) & 0xff) << 48) + ((int64(array[2]) & 0xff) << 40) + ((int64(array[3]) & 0xff) << 32) + ((int64(array[4]) & 0xff) << 24) + ((int64(array[5]) & 0xff) << 16) + ((int64(array[6]) & 0xff) << 8) + ((int64(array[7]) & 0xff) << 0)
}

func BytesToInt64A(array []byte, offset int) int64 {
	return ((int64(array[0+offset]) & 0xff) << 56) + ((int64(array[1+offset]) & 0xff) << 48) + ((int64(array[2+offset]) & 0xff) << 40) + ((int64(array[3+offset]) & 0xff) << 32) + ((int64(array[4+offset]) & 0xff) << 24) + ((int64(array[5+offset]) & 0xff) << 16) + ((int64(array[6+offset]) & 0xff) << 8) + ((int64(array[7+offset]) & 0xff) << 0)
}

func Int32ToBytes(n int32) [4]byte {
	b := [4]byte{}
	b[3] = (byte)(n & 0xff)
	b[2] = (byte)(n >> 8 & 0xff)
	b[1] = (byte)(n >> 16 & 0xff)
	b[0] = (byte)(n >> 24 & 0xff)
	return b
}

func Int32ToBytesA(n int32, array []byte, offset int) {
	array[3+offset] = (byte)(n & 0xff)
	array[2+offset] = (byte)(n >> 8 & 0xff)
	array[1+offset] = (byte)(n >> 16 & 0xff)
	array[offset] = (byte)(n >> 24 & 0xff)
}

func BytesToInt32(b []byte) int32 {
	return int32(b[3]&0xff) + int32((b[2]&0xff)<<8) + int32((b[1]&0xff)<<16) + int32((b[0]&0xff)<<24)
}

func BytesToInt32A(b []byte, offset int) int32 {
	return int32(b[offset+3]&0xff) + int32((b[offset+2]&0xff))<<8 + int32((b[offset+1]&0xff)<<16) + int32((b[offset]&0xff)<<24)
}

func Uint32ToBytes(n uint32) [4]byte {
	b := [4]byte{}
	b[3] = (byte)(n & 0xff)
	b[2] = (byte)(n >> 8 & 0xff)
	b[1] = (byte)(n >> 16 & 0xff)
	b[0] = (byte)(n >> 24 & 0xff)
	return b
}

func Uint32ToBytesA(n uint32, array []byte, offset int) {
	array[3+offset] = (byte)(n)
	array[2+offset] = (byte)(n >> 8 & 0xff)
	array[1+offset] = (byte)(n >> 16 & 0xff)
	array[offset] = (byte)(n >> 24 & 0xff)
}

func BytesToUint32(array []byte) uint32 {
	return ((uint32)(array[3] & 0xff)) + ((uint32)(array[2]&0xff))<<8 + ((uint32)(array[1]&0xff))<<16 + ((uint32)(array[0]&0xff))<<24
}

func BytesToUint32A(array []byte, offset int) uint32 {
	return ((uint32)(array[offset+3] & 0xff)) + ((uint32)(array[offset+2]&0xff))<<8 + ((uint32)(array[offset+1]&0xff))<<16 + ((uint32)(array[offset]&0xff))<<24
}

func Int24ToBytes(n int32) [3]byte {
	b := [3]byte{}
	b[2] = (byte)(n & 0xff)
	b[1] = (byte)(n >> 8 & 0xff)
	b[0] = (byte)(n >> 16 & 0xff)
	return b
}

func Int24ToBytesA(n int32, array []byte, offset int) {
	array[2+offset] = (byte)(n)
	array[1+offset] = (byte)(n >> 8 & 0xff)
	array[offset] = (byte)(n >> 16 & 0xff)
}

func BytesToInt24(b []byte) int32 {
	return int32(b[2]&0xff) + int32((b[1]&0xff)<<8) + int32((b[0]&0xff)<<16)
}

func BytesToInt24A(b []byte, offset int) int32 {
	return int32(b[offset+2]&0xff) + int32((b[offset+1]&0xff))<<8 + int32((b[offset]&0xff)<<16)
}

func Uint24ToBytes(n uint32) [3]byte {
	b := [3]byte{}
	b[2] = (byte)(n & 0xff)
	b[1] = (byte)(n >> 8 & 0xff)
	b[0] = (byte)(n >> 16 & 0xff)
	return b
}

func Uint24ToBytesA(n uint32, array []byte, offset int) {
	array[2+offset] = (byte)(n)
	array[1+offset] = (byte)(n >> 8 & 0xff)
	array[offset] = (byte)(n >> 16 & 0xff)
}

func BytesToUint24(array []byte) uint32 {
	return ((uint32)(array[2] & 0xff)) + ((uint32)(array[1]&0xff))<<8 + ((uint32)(array[0]&0xff))<<16
}

func BytesToUint24A(array []byte, offset int) uint32 {
	return ((uint32)(array[offset+2] & 0xff)) + ((uint32)(array[offset+1]&0xff))<<8 + ((uint32)(array[offset]&0xff))<<16
}

func Int16ToBytes(n int16) [2]byte {
	b := [2]byte{}
	b[1] = (byte)(n & 0xff)
	b[0] = (byte)(n >> 8 & 0xff)
	return b
}

func Int16ToBytesA(n int16, array []byte, offset int) {
	array[1+offset] = (byte)(n)
	array[offset] = (byte)(n >> 8 & 0xff)
}

func BytesToInt16(b []byte) int16 {
	return int16(b[1]&0xff) + int16((b[0]&0xff)<<8)
}

func BytesToInt16A(b []byte, offset int) int16 {
	return int16(b[offset+1]&0xff) + int16((b[offset]&0xff))<<8
}

func Uint16ToBytes(n uint16) [2]byte {
	b := [2]byte{}
	b[1] = (byte)(n & 0xff)
	b[0] = (byte)(n >> 8 & 0xff)
	return b
}

func Uint16ToBytesA(n uint16, array []byte, offset int) {
	array[1+offset] = (byte)(n)
	array[offset] = (byte)(n >> 8 & 0xff)
}

func BytesToUint16(array []byte) uint16 {
	return ((uint16)(array[1] & 0xff)) + ((uint16)(array[0]&0xff))<<8
}

func BytesToUint16A(array []byte, offset int) uint16 {
	return ((uint16)(array[offset+1] & 0xff)) + ((uint16)(array[offset]&0xff))<<8
}
