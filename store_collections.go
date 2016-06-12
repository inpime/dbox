package dbox

var registredStores = make(map[string]Store)

func GetStore(name string) (Store, error) {
	store, exists := registredStores[name]

	if !exists {
		return nil, ErrNotFound
	}

	return store, nil
}
