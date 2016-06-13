package dbox

var _ Object = (*MapObject)(nil)

func NewMapObject(store Store) *MapObject {
	return &MapObject{
		meta:  make(map[string]interface{}),
		store: store,
	}
}

type MapObject struct {
	object

	meta map[string]interface{}

	store Store
}

func (f MapObject) Map() map[string]interface{} {
	return f.meta
}

// Encoder

func (f *MapObject) Encode() (err error) {
	b, err := encode(f.meta)
	f.Write(b)

	return err
}

// Decoder

func (f MapObject) Decode() error {

	return decode(&f.meta, f.Bytes())
}

func (f *MapObject) Sync() error {

	if f.IsNew() {
		f.id = NewUUID()
	}

	f.Encode()

	return f.store.Save(f)
}
