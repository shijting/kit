package encoding

import (
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
func StructToBytes(data interface{}) []byte {
	dataSize := int(reflect.TypeOf(data).Size())
	ptr := unsafe.Pointer(reflect.ValueOf(data).Pointer())
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(ptr),
		Len:  dataSize,
		Cap:  dataSize,
	}))
}
