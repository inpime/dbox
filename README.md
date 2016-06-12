# dbox
abstraction layer for store any objects to use same api

[![Build Status](https://travis-ci.org/inpime/dbox.svg?branch=master)](https://travis-ci.org/inpime/dbox)

Status. In pursuit of the ideal architecture and internal api

# Object

Any object implements the interface. 

* `ID() string` - the object identifier
* `SetID(string)`
* `IsNew() bool`
* `Bytes() []byte` - the data of the object
* `Write([]byte)`
* `Decode() error` - decode the object. Object data are decoded
* `Encode() error` - encode the object, updated data of the object
* `Sync() error` - to update data in the store

# Helper objects

* `MapObject` - object with structured data
* `RawObject` - object with raw data
* `RefObject` - a special object to store links files (e.g. a reference file by file name)

# File

Implements an `MapObject`.

* `Delete() error` - remove file in the store
* `RawData() Object` - raw data file
* `Meta() *typed.Typed` - meta data file, `typed.Typed` this is an `map[string]interface{}` with helper functions
* `MapData() *typed.Typed` - structured data file
* `Name() string` - file name
* `SetName() string` - file name
* `CreatedAt() time.Time` - file creation date
* `UpdatedAt() time.Time` - file creation date
* `SetMapDataStore(Store)` - set a repository for structured data file
* `SetRawDataStore(Store)` - set a repository for raw data file

The best guide is the test.

# Storages

- [x] Memory. Stored in the memory (e.g. for testing)
- [x] The local file system. Need to set path store 
- [ ] amazon s3
- [ ] google storage

Notes. Library [typed](gebv/typed) is temporary.
