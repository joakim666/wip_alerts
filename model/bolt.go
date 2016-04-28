package model

import (
	"fmt"
	"reflect"

	"github.com/boltdb/bolt"
	"github.com/golang/glog"
)

type ChildID string
type ParentID string

func BoltSaveObject(bucket *bolt.Bucket, key string, obj interface{}) error {
	bytes, err := serialize(obj)
	if err != nil {
		return err
	}

	err = bucket.Put([]byte(key), bytes)
	if err != nil {
		return fmt.Errorf("Failed to save object: %s", err)
	}

	return nil
}

// TODO rename to SaveChildObjects
func BoltSaveAccountObjects(db *bolt.DB, accountUUID ParentID, bucketName string, objs *map[string]PersistanceID) error {
	glog.Infof("Saving %s for account %s", bucketName, accountUUID)

	err := db.Update(func(tx *bolt.Tx) error {
		mb := tx.Bucket([]byte(bucketName)) // main bucket

		nb, err := mb.CreateBucketIfNotExists([]byte(accountUUID)) // nested bucket
		if err != nil {
			return fmt.Errorf("Failed to create nested %s bucket for account %s: %s", bucketName, accountUUID, err)
		}

		for _, v := range *objs {
			glog.Infof("Saving object %s", v.PersistanceID())
			err := BoltSaveObject(nb, v.PersistanceID(), v)
			if err != nil {
				return fmt.Errorf("Failed to save object: %s", err)
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("Failed to save %s for account %s: %s", bucketName, accountUUID, err)
	}

	return nil
}

// TODO rename to GetChildObjects
func BoltGetAccountObjects(db *bolt.DB, accountUUID ParentID, bucketName string, t reflect.Type) (*map[string]PersistanceID, error) {
	objs := make(map[string]PersistanceID)

	err := db.View(func(tx *bolt.Tx) error {
		mb := tx.Bucket([]byte(bucketName))  // main bucket
		nb := mb.Bucket([]byte(accountUUID)) // nested bucket
		if nb == nil {
			// no nested bucket => no objects => return nil
			glog.Infof("Account %s is missing nested bucket for %s", accountUUID, bucketName)
			return nil
		}

		err := nb.ForEach(func(k, v []byte) error {
			o := reflect.New(t).Interface() // make new instance to deserialize into
			err := deserialize(&v, o)
			if err != nil {
				return fmt.Errorf("Failed to deserialize object: %s", err)
			}

			p, _ := reflect.ValueOf(o).Interface().(PersistanceID) // cast to PersistanceID to save in map to return

			objs[p.PersistanceID()] = p

			return nil
		})

		return err
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to get %s objects: %s", bucketName, err)
	}

	return &objs, nil
}

// BoltGetObject returns the object with the given 'objID' and the accountID it belongs to or nil if none is found
func BoltGetObject(db *bolt.DB, bucketName string, objID string, t reflect.Type) (*PersistanceID, *ParentID, error) {
	var obj *PersistanceID
	var accountID ParentID // TODO rename

	err := db.View(func(tx *bolt.Tx) error {
		mb := tx.Bucket([]byte(bucketName)) // main bucket

		err := mb.ForEach(func(k, v []byte) error {
			if v == nil {
				// nested bucket
				nb := mb.Bucket(k) // nested bucket
				if nb == nil {
					return fmt.Errorf("Failed to open nested bucket")
				}
				err := nb.ForEach(func(kk, vv []byte) error {
					o := reflect.New(t).Interface() // make new instance to deserialize into
					err := deserialize(&vv, o)
					if err != nil {
						return fmt.Errorf("Failed to deserialize object: %s", err)
					}

					p, _ := reflect.ValueOf(o).Interface().(PersistanceID) // cast to PersistanceID to save in map to return

					if p.PersistanceID() == objID {
						obj = &p
						accountID = ParentID(string(k))
					}

					return nil
				})
				return err

			}
			return nil
		})

		return err
	})
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to get %s object: %s", bucketName, err)
	}

	return obj, &accountID, nil
}

func BoltGetObjects(db *bolt.DB, bucketName string, t reflect.Type) (*map[string][]PersistanceID, error) {
	objs := make(map[string][]PersistanceID)

	err := db.View(func(tx *bolt.Tx) error {
		mb := tx.Bucket([]byte(bucketName)) // main bucket

		err := mb.ForEach(func(k, v []byte) error {
			if v == nil {
				// nested bucket
				nb := mb.Bucket(k) // nested bucket
				if nb == nil {
					return fmt.Errorf("Failed to open nested bucket")
				}
				err := nb.ForEach(func(k, v []byte) error {
					o := reflect.New(t).Interface() // make new instance to deserialize into
					err := deserialize(&v, o)
					if err != nil {
						return fmt.Errorf("Failed to deserialize object: %s", err)
					}

					p, _ := reflect.ValueOf(o).Interface().(PersistanceID) // cast to PersistanceID to save in map to return

					objs[p.PersistanceID()] = append(objs[p.PersistanceID()], p)

					return nil
				})
				return err

			}
			return nil
		})

		return err
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to get %s objects: %s", bucketName, err)
	}

	return &objs, nil
}

func BoltMap(objs interface{}) *map[string]PersistanceID {
	res := make(map[string]PersistanceID)

	// unwrap pointer to actual map
	r := reflect.ValueOf(objs)
	m := r.Elem()

	// then iterate over the map and cast values to PersistanceID
	for _, key := range m.MapKeys() {
		v := m.MapIndex(key)
		p, _ := v.Interface().(PersistanceID)
		res[p.PersistanceID()] = p
	}
	return &res
}

func BoltSingle(obj interface{}) *map[string]PersistanceID {
	res := make(map[string]PersistanceID)

	// unwrap pointer to actual object
	r := reflect.ValueOf(obj)
	p, _ := r.Elem().Interface().(PersistanceID)
	res[p.PersistanceID()] = p

	return &res
}
