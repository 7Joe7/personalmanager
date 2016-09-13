package db

import (
	"github.com/boltdb/bolt"
	"fmt"
)

func PrintoutDbContents(path string) {
	err := db.View(func (tx *bolt.Tx) error {
		return tx.ForEach(func (bucketName []byte, b *bolt.Bucket) error {
			fmt.Printf("%s contains:\n", string(bucketName))
			return b.ForEach(func (k, v []byte) error {
				fmt.Printf("'%s': '%s'", string(k), string(v))
				return nil
			})
		})
	})
	if err != nil {
		panic(err)
	}
}
