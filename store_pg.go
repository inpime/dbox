package dbox

import (
	"fmt"

	pg "gopkg.in/pg.v4"
)

var dbschema = []string{
	`CREATE TABLE IF NOT EXISTS %s (
            id text NOT NULL PRIMARY KEY, 
            data bytea
            )`,
}

func NewPGStore(db *pg.DB, name string) (*PGStore, error) {
	store := &PGStore{
		tname: name,
		db:    db,
	}

	if err := executeQueries(db, name, dbschema); err != nil {
		return nil, err
	}

	return store, nil
}

type PGStore struct {
	tname string
	db    *pg.DB
}

func (s PGStore) Get(id string, obj Object) error {
	var (
		data []byte
	)
	sql := fmt.Sprintf(`SELECT data FROM %s WHERE id = ?`, s.tname)
	_, err := s.db.QueryOne(pg.Scan(&data), sql, id)

	if err != nil {
		if err == pg.ErrNoRows {
			return ErrNotFound
		}

		return err
	}

	obj.Write(superDecodeData(data))
	obj.SetID(id)

	return obj.Decode()
}

func (s *PGStore) save(obj Object) error {
	var (
		id   = obj.ID()
		data = superEncodeData(obj.Bytes())
	)

	sql := fmt.Sprintf(`INSERT INTO %s (id, data)
		VALUES (?, ?)
		ON CONFLICT (id)
		DO UPDATE SET
		data=?`, s.tname)

	_, err := s.db.Exec(sql,
		id,
		data,
		data)

	return err
}

func (s *PGStore) saveFileRefs(file *File) error {
	ref := NewRefObject(s)
	ref.SetID(file.Name())
	ref.SetRefID(file.ID())

	return s.save(ref)
}

func (s *PGStore) Save(obj Object) (err error) {
	if err = s.save(obj); err != nil {
		return
	}

	switch obj := obj.(type) {
	case *File:
		err = s.saveFileRefs(obj)
	}

	return
}

func (s *PGStore) delete(obj Object) (err error) {
	var (
		id = obj.ID()
	)
	sql := fmt.Sprintf(`DELETE FROM %s WHERE id = ?;`, s.tname)

	_, err = s.db.Exec(sql, id)

	return
}

func (s *PGStore) Delete(obj Object) (err error) {
	if err = s.delete(obj); err != nil {
		return
	}

	switch obj := obj.(type) {
	case *File:
		// load refs file
		fileRef := NewRefObject(s)

		if err = s.Get(obj.Name(), fileRef); err != nil {
			return
		}

		err = s.delete(fileRef)
	}

	return
}

func (s PGStore) Type() StoreType {
	return PGStoreType
}

// GetByName Get object by name (via ReObject)
func (s PGStore) GetByName(name string, obj Object) error {
	ref := NewRefObject(&s)

	if err := s.Get(name, ref); err != nil {
		return err
	}

	return s.Get(ref.RefID(), obj)
}

// helper

func executeQueries(db *pg.DB, tname string, queries []string) error {
	for _, q := range queries {
		_, err := db.Exec(fmt.Sprintf(q, tname))

		if err != nil {
			return err
		}
	}

	return nil
}
