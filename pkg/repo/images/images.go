package images

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"strings"

	"github.com/dgraph-io/badger/v3"
	utils "github.com/fredbi/demo-api/pkg/images-utils"
	"github.com/fredbi/demo-api/pkg/repo"
)

const (
	thumbsPrefix = "thumbs/"
	thumbW       = 100
	thumbH       = 100
)

type imagesRepo struct {
	db *badger.DB
}

// New instance of a repository for images.
//
// NOTE: for the sake of this demo implementation, the KV store is backed by memory.
// It is trivial to back this store by actual files, and this would require more
// configuration options (where to store, etc).
//
// In order to maximize the retrieval of lists with thumbnails, we maintain a separate set of keys
// with thumbnails only. Iterating through this set of keys is assumed to remain fast enough, the
// full image fetch is assumed to be on one single image.
//
// The choice of the prefix includes a path separator, assuming the file name is a base name only.
func New() repo.ImagesRepo {
	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	if err != nil {
		log.Fatal(err)
	}

	return &imagesRepo{
		db: db,
	}
}

func (r *imagesRepo) Get(key string) (io.ReadCloser, error) {
	log.Printf("DEBUG: Get[%s]", key)
	var buf []byte
	err := r.db.View(func(txn *badger.Txn) error {
		item, e := txn.Get([]byte(key))
		if e != nil {
			if errors.Is(e, badger.ErrKeyNotFound) {
				return repo.ErrNotFound
			}

			return e
		}

		buf, e = item.ValueCopy(nil)
		return e
	})
	if err != nil {
		return nil, err
	}

	// NOTE: in this crude implementation, we buffer the writes to the persistent
	// store. Our interface allows for streaming I/Os to the database, if the driver
	// allows this (we can do that for instance, with a more involved badgerdb implementation,
	// or alternatively with postgres large objects).
	return io.NopCloser(bytes.NewReader(buf)), nil
}

func (r *imagesRepo) List() ([]repo.Image, error) {
	log.Printf("DEBUG: List")
	result := make([]repo.Image, 0, 10) // preallocate only 10 entries

	err := r.db.View(func(txn *badger.Txn) error {
		iterator := txn.NewIterator(badger.IteratorOptions{
			PrefetchSize:   10,
			PrefetchValues: true,
			Prefix:         []byte(thumbsPrefix),
		})

		defer iterator.Close()
		for iterator.Rewind(); iterator.Valid(); iterator.Next() {
			item := iterator.Item()
			thumb, e := item.ValueCopy(nil)
			if e != nil {
				return e
			}

			encoded := base64.StdEncoding.EncodeToString(thumb)
			result = append(result, repo.Image{
				Key:   strings.TrimPrefix(string(item.Key()), thumbsPrefix),
				Thumb: encoded,
			})
		}

		return nil

	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *imagesRepo) Create(key string, rdr io.ReadCloser) error {
	log.Printf("DEBUG: Create[%s]", key)
	err := r.db.Update(func(txn *badger.Txn) error {
		_, e := txn.Get([]byte(key))
		if e == nil {
			return errors.New("image already exists")
		}

		if !errors.Is(e, badger.ErrKeyNotFound) {
			return e
		}

		buf, e := ioutil.ReadAll(rdr)
		if e != nil {
			return e
		}

		rdr := bytes.NewReader(buf)
		thumb, e := utils.Resize(rdr, thumbW, thumbH)
		if e != nil {
			return e
		}

		e = txn.Set([]byte(key), buf)
		if e != nil {
			return e
		}

		e = txn.Set([]byte(thumbsPrefix+key), thumb)

		return e
	})

	return err
}

func (r *imagesRepo) Update(key string, rdr io.ReadCloser) error {
	log.Printf("DEBUG: Update[%s]", key)
	err := r.db.Update(func(txn *badger.Txn) error {
		_, e := txn.Get([]byte(key))
		if e != nil {
			if errors.Is(e, badger.ErrKeyNotFound) {
				return repo.ErrNotFound
			}

			return e
		}

		buf, e := ioutil.ReadAll(rdr)
		if e != nil {
			return e
		}

		rdr := bytes.NewReader(buf)
		thumb, e := utils.Resize(rdr, thumbW, thumbH)
		if e != nil {
			return e
		}

		e = txn.Set([]byte(key), buf)
		if e != nil {
			return e
		}

		e = txn.Set([]byte(thumbsPrefix+key), thumb)

		return e
	})

	return err
}

func (r *imagesRepo) Delete(key string) error {
	log.Printf("DEBUG: Delete[%s]", key)
	err := r.db.Update(func(txn *badger.Txn) error {
		e := txn.Delete([]byte(key))
		if e != nil && errors.Is(e, badger.ErrKeyNotFound) {
			return repo.ErrNotFound
		}

		e = txn.Delete([]byte(thumbsPrefix + key))
		if e != nil && errors.Is(e, badger.ErrKeyNotFound) {
			return repo.ErrNotFound
		}

		return e
	})

	return err
}
