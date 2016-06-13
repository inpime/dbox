# dbox
abstraction layer for store any objects to use same api

[![Build Status](https://travis-ci.org/inpime/dbox.svg?branch=master)](https://travis-ci.org/inpime/dbox)
[![Go Report Card](https://goreportcard.com/badge/github.com/inpime/dbox)](https://goreportcard.com/report/github.com/inpime/dbox)

Status. In pursuit of the ideal architecture and internal api.

TODO: example of using

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
Any file refers to bucket. Bucket defines how file storage (sets stores for internal objects of file).

* `Delete() error` - remove file in the store
* `RawData() Object` - raw data file
* `Meta() map[string]interface{}` - meta data file
* `MapData() map[string]interface{}` - structured data file
* `Name() string` - file name
* `SetName() string` - file name
* `Bucket()` - bucket name (bucket must exist)
* `EntityType() EntityType` - type of file (for high-level logic). 
* `CreatedAt() time.Time` - file creation date
* `UpdatedAt() time.Time` - file creation date
* `SetMapDataStore(Store)` - set a repository for structured data file
* `SetRawDataStore(Store)` - set a repository for raw data file

The best guide is the test.

TODO: Simple examples for quick start.

# Storages

- [x] Memory. Stored in the memory (e.g. for testing)
- [x] The local file system. Need to set path store 
- [ ] [amazon s3](https://aws.amazon.com/s3)
- [ ] [boltdb](https://github.com/boltdb/bolt)
- [ ] [google storage](https://cloud.google.com/storage/)