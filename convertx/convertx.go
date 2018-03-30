package convertx

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Accessible to external packages /////////////////////////////////////////////

// Atoi ...
func Atoi(str interface{}) int {
	return atoi(str)
}

// Atoui ...
func Atoui(str interface{}) uint {
	return atoui(str)
}

// ToInt64 should convert given value to int64.
func ToInt64(v interface{}) (d int64, err error) {
	return toInt64(v)
}

// BytesToInt64 should convert bytes to int64.
func BytesToInt64(buf []byte) int64 {
	return bytesToInt64(buf)
}

// Int64ToBytes should convert int64 to bytes.
func Int64ToBytes(i int64) []byte {
	return int64ToBytes(i)
}

// Underlying functions ////////////////////////////////////////////////////////

func atoi(str interface{}) (i int) {
	if str == nil {
		return 0
	}
	i, _ = strconv.Atoi(strings.Trim(str.(string), " "))
	return
}

func atoui(str interface{}) uint {
	if str == nil {
		return 0
	}
	u, _ := strconv.Atoi(strings.Trim(str.(string), " "))
	return uint(u)
}

func toInt64(v interface{}) (d int64, err error) {
	val := reflect.ValueOf(v)
	switch v.(type) {
	case int, int8, int16, int32, int64:
		d = val.Int()
	case uint, uint8, uint16, uint32, uint64:
		d = int64(val.Uint())
	default:
		err = fmt.Errorf("ToInt64 need numeric not `%T`", v)
	}
	return
}

func bytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

func int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

////////////////////////////////////////////////////////////////////////////////
