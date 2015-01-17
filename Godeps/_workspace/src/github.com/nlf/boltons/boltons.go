package boltons

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"

	"code.google.com/p/go-uuid/uuid"

	"github.com/boltdb/bolt"
)

type DB struct {
	bolt *bolt.DB
}

func Open(path string, mode os.FileMode, options *bolt.Options) (*DB, error) {
	db, err := bolt.Open(path, mode, options)
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

type parsedBucket struct {
	name   []byte
	values map[string]reflect.Value
	fields map[string]reflect.StructField
}

func parseInput(s interface{}, writable bool) (parsedBucket, error) {
	errMsg := errors.New("Expected struct pointer")
	bucket := parsedBucket{}

	sType := reflect.TypeOf(s)
	if sType.Kind() == reflect.Ptr {
		sType = sType.Elem()
	} else {
		if writable {
			return bucket, errMsg
		}
	}

	sValue := reflect.Indirect(reflect.ValueOf(s))
	if sValue.Kind() != reflect.Struct {
		return bucket, errMsg
	}

	bucket.name = []byte(sType.Name())
	bucket.values = make(map[string]reflect.Value)
	bucket.fields = make(map[string]reflect.StructField)

	for i := 0; i < sValue.NumField(); i++ {
		fValue := sValue.Field(i)
		fType := sType.Field(i)

		bucket.values[fType.Name] = fValue
		bucket.fields[fType.Name] = fType
	}

	return bucket, nil
}

func (db *DB) Save(s interface{}) error {
	bucket, err := parseInput(s, true)
	if err != nil {
		return err
	}

	err = db.bolt.Update(func(tx *bolt.Tx) error {
		outer, err := tx.CreateBucketIfNotExists(bucket.name)
		if err != nil {
			return err
		}

		id := bucket.values["ID"]
		if id.String() == "" {
			id.SetString(uuid.New())
		}

		inner, err := outer.CreateBucketIfNotExists([]byte(id.String()))
		if err != nil {
			return err
		}

		for key, value := range bucket.values {
			bVal, err := json.Marshal(value.Interface())
			if err != nil {
				return nil
			}

			err = inner.Put([]byte(key), bVal)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func (db *DB) Get(s interface{}) error {
	bucket, err := parseInput(s, true)
	if err != nil {
		return err
	}

	err = db.bolt.View(func(tx *bolt.Tx) error {
		outer := tx.Bucket(bucket.name)

		id := bucket.values["ID"]
		if id.String() == "" {
			return errors.New("Unable to fetch without an ID")
		}

		inner := outer.Bucket([]byte(id.String()))

		for key, value := range bucket.values {
			bVal := inner.Get([]byte(key))

			out := reflect.New(value.Type()).Interface()
			err := json.Unmarshal(bVal, &out)
			if err != nil {
				return err
			}

			if out != nil {
				value.Set(reflect.Indirect(reflect.ValueOf(out)))
			}
		}

		return nil
	})

	return err
}

func (db *DB) Update(s interface{}, changes map[string]interface{}) error {
	bucket, err := parseInput(s, true)
	if err != nil {
		return err
	}

	err = db.bolt.Update(func(tx *bolt.Tx) error {
		outer := tx.Bucket(bucket.name)

		id := bucket.values["ID"]
		if id.String() == "" {
			return errors.New("Unable to fetch without an ID")
		}

		inner := outer.Bucket([]byte(id.String()))

		for key, value := range bucket.values {
			val, ok := changes[key]
			if !ok {
				jsonKey := bucket.fields[key].Tag.Get("json")
				val, ok = changes[jsonKey]
			}
			if ok {
				bVal, err := json.Marshal(val)
				if err != nil {
					return nil
				}

				err = inner.Put([]byte(key), bVal)
				if err != nil {
					return err
				}

				value.Set(reflect.Indirect(reflect.ValueOf(val)))
			} else {
				bVal := inner.Get([]byte(key))

				out := reflect.New(value.Type()).Interface()
				err := json.Unmarshal(bVal, &out)
				if err != nil {
					return err
				}

				if out != nil {
					value.Set(reflect.Indirect(reflect.ValueOf(out)))
				}
			}
		}

		return nil
	})

	return err
}

func (db *DB) All(s interface{}) error {
	errMsg := errors.New("Expected pointer to struct slice")

	if reflect.TypeOf(s).Kind() != reflect.Ptr {
		return errMsg
	}

	sValue := reflect.Indirect(reflect.ValueOf(s))
	if sValue.Kind() != reflect.Slice {
		return errMsg
	}

	sType := sValue.Type().Elem()
	if sType.Kind() != reflect.Struct {
		return errMsg
	}

	bucketName := sType.Name()
	err := db.bolt.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return nil
		}

		cursor := bucket.Cursor()

		for key, _ := cursor.First(); key != nil; key, _ = cursor.Next() {
			inner := bucket.Bucket(key)
			innerCursor := inner.Cursor()

			member := reflect.New(sType).Elem()
			for key, value := innerCursor.First(); key != nil; key, value = innerCursor.Next() {
				field := member.FieldByName(string(key))
				out := reflect.New(field.Type()).Interface()

				err := json.Unmarshal(value, &out)
				if err != nil {
					return err
				}

				if out != nil {
					field.Set(reflect.Indirect(reflect.ValueOf(out)))
				}
			}

			sValue.Set(reflect.Append(sValue, member))
		}

		return nil
	})

	return err
}

func (db *DB) Keys(s interface{}) ([]string, error) {
	keys := []string{}
	bucket, err := parseInput(s, false)
	if err != nil {
		return keys, err
	}

	err = db.bolt.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucket.name)
		if bucket == nil {
			return nil
		}
		cursor := bucket.Cursor()

		for key, _ := cursor.First(); key != nil; key, _ = cursor.Next() {
			keys = append(keys, string(key))
		}

		return nil
	})

	return keys, err
}

func (db *DB) Exists(s interface{}) (bool, error) {
	exists := false
	bucket, err := parseInput(s, false)
	if err != nil {
		return exists, err
	}

	err = db.bolt.View(func(tx *bolt.Tx) error {
		outer := tx.Bucket(bucket.name)
		id := bucket.values["ID"].String()
		if id == "" {
			return nil
		}

		inner := outer.Bucket([]byte(id))
		exists = inner != nil
		return nil
	})

	return exists, err
}

func (db *DB) Delete(s interface{}) error {
	bucket, err := parseInput(s, false)
	if err != nil {
		return err
	}

	err = db.bolt.Update(func(tx *bolt.Tx) error {
		outer := tx.Bucket(bucket.name)
		id := bucket.values["ID"].String()
		if id == "" {
			return nil
		}

		err := outer.DeleteBucket([]byte(id))
		return err
	})

	return err
}

func (db *DB) Close() {
	db.bolt.Close()
}
