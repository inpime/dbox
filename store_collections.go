package dbox

var registredStores = make(map[string]Store)

func RegistryStore(name string, store Store) {
	registredStores[name] = store
}

func MustStore(name string) Store {
	store, exists := registredStores[name]

	if !exists {
		return nil
	}

	return store
}
