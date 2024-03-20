package mapstructure

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"uniswap-transaction/pkg/utils/xtime"
)

type Basic struct {
	Vstring     string
	Vint        int
	Vint8       int8
	Vint16      int16
	Vint32      int32
	Vint64      int64
	Vuint       uint
	Vbool       bool
	Vfloat      float64
	Vextra      string
	vsilent     bool
	Vdata       interface{}
	VjsonInt    int
	VjsonUint   uint
	VjsonUint64 uint64
	VjsonFloat  float64
	//VjsonNumber json.Number
}

type EmbeddedAndNamed struct {
	Basic
	Named   Basic
	Vunique string
}

type TTime struct {
	Id   int       `json:"id"`
	Time time.Time `json:"time"`
}

type TXTime struct {
	Id   int        `json:"id"`
	Time xtime.Time `json:"time"`
}

type TTimePointer struct {
	Id   int        `json:"id"`
	Time *time.Time `json:"time"`
}

type TXTimePointer struct {
	Id   int         `json:"id"`
	Time *xtime.Time `json:"time"`
}

func TestDecode(t *testing.T) {
	input := EmbeddedAndNamed{
		Basic:   Basic{Vstring: "foo"},
		Named:   Basic{Vstring: "baz"},
		Vunique: "bar",
	}
	var m map[string]interface{}
	err := Decode(input, &m)
	assert.NoError(t, err)
	fmt.Println(m)
}

func TestDecodeTime(t *testing.T) {
	m := map[string]interface{}{
		"id":   1,
		"time": "2022-06-10 08:00:00",
	}
	var tTime TTime
	err := Decode(m, &tTime)
	fmt.Println(tTime)
	assert.NoError(t, err)
}

func TestDecodeXTime(t *testing.T) {
	m := map[string]interface{}{
		"id":   1,
		"time": "2022-06-10 08:00:00",
	}
	var tXTime TXTime
	err := Decode(m, &tXTime)
	fmt.Println(tXTime)
	assert.NoError(t, err)
}

func TestDecodeTimePointer(t *testing.T) {
	m := map[string]interface{}{
		"id":   1,
		"time": nil,
	}
	var tTime TTimePointer
	err := Decode(m, &tTime)
	fmt.Println(tTime)
	assert.NoError(t, err)
}

func TestDecodeXTimePointer(t *testing.T) {
	m := map[string]interface{}{
		"id":   1,
		"time": nil,
	}
	var tXTime TXTimePointer
	err := Decode(m, &tXTime)
	fmt.Println(tXTime)
	assert.NoError(t, err)
}
