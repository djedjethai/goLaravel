package cache

import (
	"github.com/dgraph-io/badger/v3"
	"time"
)

type BadgerCache struct {
	Conn   *badger.DB
	Prefix string
}

func (b *BadgerCache) Has(str string) (bool, error) {
	_, err := b.Get(str)
	if err != nil {
		return false, nil
	}

	return true, nil
}

func (b *BadgerCache) Get(str string) (interface{}, error) {
	var fromCache []byte
	err := b.Conn.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(str))
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			fromCache = append([]byte{}, val...)
			return nil
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	decoded, err := decode(string(fromCache))
	if err != nil {
		return nil, err
	}

	item := decoded[str]

	return item, nil
}

func (b *BadgerCache) Set(str string, value interface{}, expires ...int) error {
	entry := Entry{}

	entry[str] = value

	encoded, err := encode(entry)
	if err != nil {
		return err
	}

	if len(expires) > 0 {
		err = b.Conn.Update(func(txn *badger.Txn) error {
			e := badger.NewEntry([]byte(str), encoded).WithTTL(time.Second * time.Duration(expires[0]))
			err = txn.SetEntry(e)
			return err
		})
	} else {
		err = b.Conn.Update(func(txn *badger.Txn) error {
			e := badger.NewEntry([]byte(str), encoded)
			err = txn.SetEntry(e)
			return err
		})

	}

	return nil
}

func (b *BadgerCache) Forget(str string) error {

	err := b.Conn.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(str))
		return err
	})

	return err
}

func (b *BadgerCache) EmptyByMatch(str string) error {
	return b.emptyByMatch(str)
}

func (b *BadgerCache) Empty() error {
	return b.emptyByMatch("")
}

func (b *BadgerCache) emptyByMatch(str string) error {
	// search the entire badger cached for the keys we have to delete
	deleteKeys := func(keysForDelete [][]byte) error {
		if err := b.Conn.Update(func(txn *badger.Txn) error {
			for _, key := range keysForDelete {
				if err := txn.Delete(key); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return err
		}
		return nil
	}
	collectSize := 100000

	err := b.Conn.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.AllVersions = false
		opts.PrefetchValues = false
		// it for iterator
		it := txn.NewIterator(opts)
		defer it.Close()

		keysForDelete := make([][]byte, 0, collectSize)
		keysCollected := 0

		// the iterator is seeking/itering
		for it.Seek([]byte(str)); it.ValidForPrefix([]byte(str)); it.Next() {
			key := it.Item().KeyCopy(nil)
			keysForDelete = append(keysForDelete, key)
			keysCollected++
			// check if the collected key do not hit the limit size(100000)
			// make sure we do delete more than 100000 items at the time
			// it would take too long
			if keysCollected == collectSize {
				if err := deleteKeys(keysForDelete); err != nil {
					return err
				}
			}
		}
		// if we pass the previous condition, we delete the actual keysForDelete
		if keysCollected > 0 {
			if err := deleteKeys(keysForDelete); err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

// func (b *BadgerCache) getKeys(pattern string) ([]string, error) {
// 	conn := c.Conn.Get()
// 	defer conn.Close()
//
// 	iter := 0
// 	keys := []string{}
//
// 	for {
// 		arr, err := redis.Values(conn.Do("SCAN", iter, "MATCH", fmt.Sprintf("%s*", pattern)))
// 		if err != nil {
// 			return keys, err
// 		}
//
// 		iter, _ = redis.Int(arr[0], nil)
// 		k, _ := redis.Strings(arr[1], nil)
// 		keys = append(keys, k...)
//
// 		if iter == 0 {
// 			break
// 		}
// 	}
//
// 	return keys, nil
//
// }
