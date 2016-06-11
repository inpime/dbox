package dbox

import "gopkg.in/vmihailenco/msgpack.v2"

func encode(v interface{}) ([]byte, error) {
	return msgpack.Marshal(v)
}

func decode(v interface{}, b []byte) error {

	return msgpack.Unmarshal(b, v)
}
