package dbox

var _ Object = (*RawObject)(nil)

func NewRawObject(store Store) *RawObject {
	return &RawObject{
		store: store,
	}
}

type RawObject struct {
	object

	store Store
}

// Encoder

func (f *RawObject) Encode() error {
	return nil
}

// Decoder

func (f RawObject) Decode() error {
	return nil
}

func (f *RawObject) Sync() error {

	if f.IsNew() {
		f.id = NewUUID()
	}

	return f.store.Save(f)
}
