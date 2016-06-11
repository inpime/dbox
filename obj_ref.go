package dbox

var _ Object = (*RefObject)(nil)

func NewRefObject(store Store) *RefObject {
	return &RefObject{
		store: store,
	}
}

// RefObject link filename to a file id
type RefObject struct {
	object

	store Store
}

// Name file name as id 
func (f RefObject) RefID() string {
    return string(f.data)
}

func (f *RefObject) SetRefID(v string) {
	f.data = []byte(v)
    
}

// ID file id as data value
func (f RefObject) ID() string {
    return f.id
}

func (f *RefObject) SetID(v string) {
    f.id = v
}

// Encoder

func (f *RefObject) Encode() error {
	return nil
}

// Decoder

func (f RefObject) Decode() error {
	return nil
}

func (f *RefObject) Sync() error {

	if f.IsNew() {
		return ErrEmptyName
	}

	return f.store.Save(f)
}
