package pgxmap

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	String string    `json:"s"`
	Number int64     `json:"n"`
	Float  float64   `json:"f"`
	Time   time.Time `json:"t"`
	Bool   bool      `json:"b"`
}

func TestAState(t *testing.T) {
	s1 := State[testStruct]{
		Data: testStruct{
			String: "hello",
			Number: time.Now().UnixNano(),
			Float:  float64(time.Now().UnixNano()) / float64(time.Second),
			Time:   time.Now(),
			Bool:   true,
		},
	}

	err := s1.Scan(nil)
	assert.Error(t, err)

	d, err := s1.Value()
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, d)
	t.Log(string(d.([]byte)))

	var s2 State[testStruct]
	if err = s2.Scan(d); err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, s2)
	assert.NotNil(t, s2.Data)

	assert.Equal(t, s1.Data.String, s2.Data.String)
	assert.Equal(t, s1.Data.Number, s2.Data.Number)
	assert.Equal(t, s1.Data.Time.UTC(), s2.Data.Time.UTC())
	assert.Equal(t, s1.Data.Bool, s2.Data.Bool)
}
