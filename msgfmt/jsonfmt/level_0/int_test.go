package test

import (
	"testing"
	"reflect"
	"unsafe"
	"github.com/stretchr/testify/require"
	"github.com/v2pro/plz/msgfmt/jsonfmt"
)

func Test_int8(t *testing.T) {
	should := require.New(t)
	encoder := jsonfmt.EncoderOf(reflect.TypeOf(int8(1)))
	should.Equal("-1", string(encoder.Encode(nil,nil, jsonfmt.PtrOf(int8(-1)))))
}

func Test_uint8(t *testing.T) {
	should := require.New(t)
	encoder := jsonfmt.EncoderOf(reflect.TypeOf(uint8(1)))
	should.Equal("222", string(encoder.Encode(nil,nil, jsonfmt.PtrOf(uint8(222)))))
}

func Test_int16(t *testing.T) {
	should := require.New(t)
	encoder := jsonfmt.EncoderOf(reflect.TypeOf(int16(1)))
	should.Equal("222", string(encoder.Encode(nil,nil, jsonfmt.PtrOf(int16(222)))))
}

func Test_uint16(t *testing.T) {
	should := require.New(t)
	encoder := jsonfmt.EncoderOf(reflect.TypeOf(uint16(1)))
	should.Equal("222", string(encoder.Encode(nil,nil, jsonfmt.PtrOf(uint16(222)))))
}

func Test_int32(t *testing.T) {
	should := require.New(t)
	encoder := jsonfmt.EncoderOf(reflect.TypeOf(int32(1)))
	should.Equal("222", string(encoder.Encode(nil,nil, jsonfmt.PtrOf(int32(222)))))
}

func Test_uint32(t *testing.T) {
	should := require.New(t)
	encoder := jsonfmt.EncoderOf(reflect.TypeOf(uint32(1)))
	should.Equal("222", string(encoder.Encode(nil,nil, jsonfmt.PtrOf(uint32(222)))))
}

func Test_int64(t *testing.T) {
	should := require.New(t)
	encoder := jsonfmt.EncoderOf(reflect.TypeOf(int64(1)))
	should.Equal("222", string(encoder.Encode(nil,nil, jsonfmt.PtrOf(int64(222)))))
}

func Test_uint64(t *testing.T) {
	should := require.New(t)
	encoder := jsonfmt.EncoderOf(reflect.TypeOf(uint64(1)))
	should.Equal("222", string(encoder.Encode(nil,nil, jsonfmt.PtrOf(uint64(222)))))
}

func Test_int(t *testing.T) {
	should := require.New(t)
	encoder := jsonfmt.EncoderOf(reflect.TypeOf(1))
	should.Equal("1", string(encoder.Encode(nil,nil, jsonfmt.PtrOf(1))))
	should.Equal("998123123", string(encoder.Encode(nil,nil, jsonfmt.PtrOf(998123123))))
	should.Equal("-998123123", string(encoder.Encode(nil,nil, jsonfmt.PtrOf(-998123123))))
}

func Test_uint(t *testing.T) {
	should := require.New(t)
	encoder := jsonfmt.EncoderOf(reflect.TypeOf(uint(1)))
	should.Equal("222", string(encoder.Encode(nil,nil, jsonfmt.PtrOf(uint(222)))))
}

func Benchmark_int(b *testing.B) {
	encoder := jsonfmt.EncoderOf(reflect.TypeOf(1))
	values := []int{998123123, 123123435, 1230}
	ptrs := make([]unsafe.Pointer, len(values))
	for i := 0; i < len(values); i++ {
		ptrs[i] = unsafe.Pointer(&values[i])
	}
	var buf []byte
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buf = buf[:0]
		buf = encoder.Encode(nil,buf, ptrs[i%3])
	}
}
