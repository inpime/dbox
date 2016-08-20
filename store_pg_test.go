package dbox

import (
	"os"
	"testing"

	"gopkg.in/pg.v4"
)

func TestPGStore_simpleStrategy(t *testing.T) {
	db := pg.Connect(
		&pg.Options{
			Addr:     os.Getenv("DB_ADDR"),
			Database: os.Getenv("DB_NAME"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASS"),
		},
	)
	defer db.Close()

	store, err := NewPGStore(db, "fortest")
	if err != nil {
		t.Fatal(err)
	}

	createSimpleStrategy(t, store, store, store)
	deleteSimpleStrategy(t, store, store, store)

}

func TestPGStore_onlyFile(t *testing.T) {
	db := pg.Connect(
		&pg.Options{
			Addr:     os.Getenv("DB_ADDR"),
			Database: os.Getenv("DB_NAME"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASS"),
		},
	)
	defer db.Close()

	store, _ := NewPGStore(db, "fortest")
	// storeMap, _ := NewPGStore(db, "fortest_s")
	// storeRaw, _ := NewPGStore(db, "fortest_r")

	file := NewFile(store)
	file.SetName("name")
	// file.SetMapDataStore(storeMap)
	// file.SetRawDataStore(storeRaw)

	mapSet(file.Meta(), "a", "b")

	file.RawData().Write([]byte("text text"))
	mapSet(file.MapData(), "map1", "v1")
	err := file.Sync()
	if err != nil {
		t.Fatal(err)
	}

	fileId := file.ID()

	// Load

	file = NewFile(store)
	// file.SetMapDataStore(storeMap)
	// file.SetRawDataStore(storeRaw)

	err = store.Get(fileId, file)
	if err != nil {
		t.Fatal(err)
	}

	err = file.Delete()
	if err != nil {
		t.Fatal(err)
	}
}
