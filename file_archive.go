package dbox

type fileArchive struct {
	Prop []byte
	Map  []byte
	Raw  []byte
}

func (f *File) Import(b []byte) error {
	archive := fileArchive{}

	if err := decode(&archive, b); err != nil {
		return err
	}

	f.MapObject.Write(archive.Prop)
	if err := f.MapObject.Decode(); err != nil {
		return err
	}

	f.mdataObj().Write(archive.Map)
	if err := f.mdataObj().Decode(); err != nil {
		return err
	}

	f.rdataObj().Write(archive.Raw)
	if err := f.rdataObj().Decode(); err != nil {
		return err
	}

	return nil
}

func (f File) Export() ([]byte, error) {

	archive := fileArchive{
		Prop: f.MapObject.Bytes(),
		Map:  f.mdataObj().Bytes(),
		Raw:  f.rdataObj().Bytes(),
	}

	return encode(archive)
}
