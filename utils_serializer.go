package dbox

import (
	"bytes"
	"gopkg.in/vmihailenco/msgpack.v2"
)

func init() {
}

func encode(v interface{}) ([]byte, error) {
	return msgpack.Marshal(v)
}

func decode(v interface{}, b []byte) error {
	dec := msgpack.NewDecoder(bytes.NewBuffer(b))
	dec.DecodeMapFunc = func(d *msgpack.Decoder) (interface{}, error) {
		n, err := d.DecodeMapLen()
		if err != nil {
			return nil, err
		}

		m := make(map[string]interface{}, n)
		for i := 0; i < n; i++ {
			mk, err := d.DecodeString()
			if err != nil {
				return nil, err
			}

			mv, err := d.DecodeInterface()
			if err != nil {
				return nil, err
			}

			m[mk] = mv
		}
		return m, nil
	}

	return dec.Decode(v)
	// return msgpack.Unmarshal(b, v)
}
