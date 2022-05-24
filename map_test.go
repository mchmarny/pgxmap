package pgxmap

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	bigInt := time.Now().UnixNano()

	m := ConfigMap{
		"int64":   bigInt,
		"int32":   int32(bigInt),
		"int16":   int16(bigInt),
		"int8":    int8(bigInt),
		"int":     int(bigInt),
		"uint64":  uint64(bigInt),
		"uint32":  uint32(bigInt),
		"uint16":  uint16(bigInt),
		"uint8":   uint8(bigInt),
		"uint":    uint(bigInt),
		"float32": float32(12345.6789),
		"float64": float64(bigInt) / float64(3),
	}

	d, err := m.Value()
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, d)

	b, err := m.Value()
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, b)
	t.Log(string(b.([]byte)))

	var m2 ConfigMap
	if err = m2.Scan(b); err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, m2)

	assert.Equal(t, m["int64"].(int64), m2["int64"].(int64))
	assert.Equal(t, m["int32"].(int32), m2["int32"].(int32))
	assert.Equal(t, m["int16"].(int16), m2["int16"].(int16))
	assert.Equal(t, m["int8"].(int8), m2["int8"].(int8))
	assert.Equal(t, m["int"].(int), m2["int"].(int))
	assert.Equal(t, m["uint64"].(uint64), m2["uint64"].(uint64))
	assert.Equal(t, m["uint32"].(uint32), m2["uint32"].(uint32))
	assert.Equal(t, m["uint16"].(uint16), m2["uint16"].(uint16))
	assert.Equal(t, m["uint8"].(uint8), m2["uint8"].(uint8))
	assert.Equal(t, m["uint"].(uint), m2["uint"].(uint))
	assert.Equal(t, m["float32"].(float32), m2["float32"].(float32))
	assert.Equal(t, m["float64"].(float64), m2["float64"].(float64))
}
