package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"reflect"
	"github.com/v2pro/plz/msgfmt/jsonfmt"
)

func Test_string(t *testing.T) {
	should := require.New(t)
	encoder := jsonfmt.EncoderOf(reflect.TypeOf(""))
	should.Equal(`"hello"`, string(encoder.Encode(nil,nil, jsonfmt.PtrOf("hello"))))
	should.Equal(`"\nhello中文"`, string(encoder.Encode(nil,nil, jsonfmt.PtrOf("\nhello中文"))))
	should.Equal(`"\nhello中文h\nello"`, string(encoder.Encode(nil,nil, jsonfmt.PtrOf("\nhello中文h\nello"))))
	should.Equal(`"\nhello中文h\nello\t"`, string(encoder.Encode(nil,nil, jsonfmt.PtrOf("\nhello中文h\nello\t"))))
}
