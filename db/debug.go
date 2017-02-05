package db

import (
	"fmt"

	"github.com/boltdb/bolt"
)

func PrintoutDbContents(path string) {
	err := db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(bucketName []byte, b *bolt.Bucket) error {
			fmt.Printf("%s contains:\n", string(bucketName))
			return b.ForEach(func(k, v []byte) error {
				fmt.Printf("'%s': '%s'\n", string(k), string(v))
				return nil
			})
		})
	})
	if err != nil {
		panic(err)
	}
}
