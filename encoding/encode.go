package encoding

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

// StringToBytes 将字符串转换为字节数组
func StringToBytes(s string) []byte {
	if s == "" {
		return []byte{}
	}
	strHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: strHeader.Data,
		Len:  strHeader.Len,
		Cap:  strHeader.Len,
	}))
}

// BytesToString 将字节数组转换为字符串
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// StructToBytes 将结构体转换为字节数组
func StructToBytes(data interface{}) ([]byte, error) {
	if data == nil {
		return []byte{}, nil
	}
	fmt.Println(reflect.TypeOf(data).Kind())
	if reflect.TypeOf(data).Kind() != reflect.Struct {
		return []byte{}, errors.New("input is not a struct")
	}
	var b []byte
	size := unsafe.Sizeof(data)
	byteHeader := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	byteHeader.Len = int(size)
	byteHeader.Cap = int(size)
	byteHeader.Data = uintptr(unsafe.Pointer(&data))
	return b, nil
}
