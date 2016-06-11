package dbox

import (
	"sync"

	"encoding/json"
)

type Object interface {
	ID() string
	SetID(string)

	Sync() error

	Bytes() []byte
	// Write(p []byte) (int, error)
	// Read(p []byte) (int, error)
	Write(p []byte)
	// Read(p []byte) (int, error)

	Decode() error
	Encode() error
}

var _ json.Marshaler = (*Objects)(nil)

type Objects struct {
	sync.RWMutex

	list     []Object
	hasNext  bool
	nextPage int
	total    int
}

func (o Objects) Len() int { return len(o.list) }

func (o Objects) HasNext() bool { return o.hasNext }

func (o Objects) NextPage() int { return o.nextPage }

func (o Objects) Total() int { return o.total }

func (o Objects) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		List     []Object
		HasNext  bool
		NextPage int
		Total    int
	}{o.list, o.hasNext, o.nextPage, o.total})
}

// object
// Abstract object
// each object must implement Sync
type object struct {
	id   string
	data []byte
}

func (o object) ID() string {
	return o.id
}

func (o *object) SetID(id string) {
	o.id = id
}

func (o object) IsNew() bool {
	return len(o.id) == 0
}

func (o object) Bytes() []byte {
	return o.data
}

func (o *object) Write(p []byte) {
	o.data = p
}
