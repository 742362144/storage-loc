package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"reflect"
)



func db() {


	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("MyBucket"))
		if err != nil {
			return fmt.Errorf("create bucket: %v", err)
		}

		if err = b.Put([]byte("answer"), []byte("42")); err != nil {
			return err
		}

		if err = b.Put([]byte("zero"), []byte("31231")); err != nil {
			return err
		}

		return nil
	})

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("MyBucket"))
		v := b.Get([]byte("noexists"))
		fmt.Println(reflect.DeepEqual(v, nil)) // false
		fmt.Println(v == nil)                  // true

		v = b.Get([]byte("zero"))
		fmt.Println(reflect.DeepEqual(v, nil)) // false
		fmt.Println(v == nil)                  // true
		fmt.Println(string(v))
		return nil
	})
}