package dbox

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// func TestMap_getterSetterAndSerializer(t *testing.T) {
// 	m := map[string]interface{}{}

// 	timeNow := time.Now()

// 	mapSet(m, "a", "b")
// 	mapSet(m, "t", timeNow)

// 	assert.Equal(t, mapString(m, "a"), "b", "get map")
// 	assert.Equal(t, mapTime(m, "t"), timeNow, "get map")

// 	//

// 	network, err := encode(m)
// 	assert.NoError(t, err, "encode map")

// 	mm := map[string]interface{}{}
// 	err = decode(&mm, network)
// 	assert.NoError(t, err, "decode map")
// 	t.Logf("%#v", mm["t"].(time.Time))

// 	assert.Equal(t, mapString(mm, "a"), "b", "get map")
// 	assert.Equal(t, mapTime(mm, "t"), timeNow, "get map")
// }

func TestMSGPack_time(t *testing.T) {
	tt := time.Now()

	network, err := encode(tt)
	assert.NoError(t, err, "encode time")

	tDecode := time.Time{}
	err = decode(&tDecode, network)
	assert.NoError(t, err, "decode time")

	assert.Equal(t, tDecode, tt, "check time")

	// ----------------
	// array times
	// ----------------

	var tArr []time.Time
	tArr = append(tArr, time.Now())
	tArr = append(tArr, time.Now())

	network, err = encode(tArr)
	assert.NoError(t, err, "encode array times")

	tArrDecode := []time.Time{}
	err = decode(&tArrDecode, network)
	assert.NoError(t, err, "decode arr times")
	assert.Equal(t, tArrDecode[0], tArr[0], "check arr times")
	assert.Equal(t, tArrDecode[1], tArr[1], "check arr times")

	// ----------------
	// map times
	// ----------------

	tMap := make(map[string]time.Time)
	tMap["a"] = time.Now()
	tMap["b"] = time.Now()

	network, err = encode(tMap)
	assert.NoError(t, err, "encode map times")

	tMapDecode := make(map[string]time.Time)

	err = decode(&tMapDecode, network)
	assert.NoError(t, err, "decode arr times")
	assert.Equal(t, tMapDecode["a"], tMap["a"], "check map times")
	assert.Equal(t, tMapDecode["b"], tMap["b"], "check map times")
}
