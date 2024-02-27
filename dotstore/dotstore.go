package dotstore

import (
	"fmt"
	"path"
	"time"

	bolt "go.etcd.io/bbolt"

	"github.com/jakenvac/DotHerd/config"
)

const (
	pen_bucket  = "pen"
	dot_bucket  = "dot"
	meta_bucket = "meta"
)

func initStore() (*bolt.DB, error) {
	dbPath := path.Join(config.DEFAULT_DOT_DIR, "dotstore.db")
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, err
	}

	db.Update(func(tx *bolt.Tx) error {
		_, penErr := tx.CreateBucketIfNotExists([]byte(pen_bucket))
		if penErr != nil {
			return penErr
		}
		_, dotErr := tx.CreateBucketIfNotExists([]byte(dot_bucket))
		if dotErr != nil {
			return dotErr
		}
		m, metaErr := tx.CreateBucketIfNotExists([]byte(meta_bucket))
		if metaErr != nil {
			return metaErr
		}

		updated := m.Get([]byte("updated"))
		if updated == nil {
			now := time.Now().Format(time.RFC3339)
			m.Put([]byte("updated"), []byte(fmt.Sprintf("%s", now)))
		}

		return nil
	})

	return db, nil
}

type DotStore struct {
	db *bolt.DB
}

func (ds *DotStore) String() string {
	return ds.db.String()
}

func (ds *DotStore) Close() error {
	return ds.db.Close()
}

func (ds *DotStore) Herd(dot, aliasInPen string) error {
	return ds.db.Update(func(tx *bolt.Tx) error {
		penBucket := tx.Bucket([]byte(pen_bucket))
		dotBucket := tx.Bucket([]byte(dot_bucket))
		metaBucket := tx.Bucket([]byte(meta_bucket))

		if err := penBucket.Put([]byte(aliasInPen), []byte(dot)); err != nil {
			return err
		}

		if err := dotBucket.Put([]byte(dot), []byte(aliasInPen)); err != nil {
			penBucket.Delete([]byte(aliasInPen))
			return err
		}

		now := time.Now().Format(time.RFC3339)
		if err := metaBucket.Put([]byte("updated"), []byte(fmt.Sprintf("%s", now))); //
		err != nil {
			penBucket.Delete([]byte(aliasInPen))
			dotBucket.Delete([]byte(dot))
		}

		return nil
	})
}

func (ds *DotStore) DotToPenAlias(dot string) (string, error) {
	var alias string
	err := ds.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dot_bucket))
		alias = string(b.Get([]byte(dot)))
		return nil
	})
	if err != nil {
		return "", err
	}

	return alias, nil
}

func (ds *DotStore) PenAliasToDot(alias string) (string, error) {
	var dot string
	err := ds.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(pen_bucket))
		dot = string(b.Get([]byte(alias)))
		return nil
	})
	if err != nil {
		return "", err
	}

	return dot, nil
}

func (ds *DotStore) Release(dot string) error {
	return ds.db.Update(func(tx *bolt.Tx) error {
		penBucket := tx.Bucket([]byte(pen_bucket))
		dotBucket := tx.Bucket([]byte(dot_bucket))
		metaBucket := tx.Bucket([]byte(meta_bucket))

		alias := string(dotBucket.Get([]byte(dot)))
		if err := penBucket.Delete([]byte(alias)); err != nil {
			return err
		}

		if err := dotBucket.Delete([]byte(dot)); err != nil {
			penBucket.Put([]byte(alias), []byte(dot))
			return err
		}

		now := time.Now().Format(time.RFC3339)
		if err := metaBucket.Put([]byte("updated"), []byte(fmt.Sprintf("%s", now))); //
		err != nil {
			penBucket.Put([]byte(alias), []byte(dot))
			dotBucket.Put([]byte(dot), []byte(alias))
		}

		return nil
	})
}

func New() (*DotStore, error) {
	db, err := initStore()
	if err != nil {
		return nil, err
	}

	return &DotStore{db: db}, nil
}
