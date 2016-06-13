package dbox

import (
	"github.com/satori/go.uuid"
	"strconv"
	"strings"
)

func NewUUID() string {
	return strings.Replace(uuid.NewV4().String(), "-", "", -1)
}

func mapSet(m map[string]interface{}, key string, v interface{}) {
	m[key] = v
}

func mapInt64(m map[string]interface{}, key string) int64 {
	value, exists := m[key]

	if exists == false {
		return 0
	}

	switch t := value.(type) {
	case int:
		return int64(t)
	case int16:
		return int64(t)
	case int32:
		return int64(t)
	case int64:
		return t
	case uint:
		return int64(t)
	case uint8:
		return int64(t)
	case uint16:
		return int64(t)
	case uint32:
		return int64(t)
	case uint64:
		return int64(t)
	case float64:
		return int64(t)
	case float32:
		return int64(t)
	case string:

		i, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			return 0
		}

		return i
	}

	return 0
}

func mapString(m map[string]interface{}, key string) string {
	value, exists := m[key]

	if exists == false {
		return ""
	}

	if n, ok := value.(string); ok {
		return n
	}
	return ""
}
