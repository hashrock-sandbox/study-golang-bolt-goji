package main

import (
    "log"
    "fmt"
    "github.com/boltdb/bolt"
    "strconv"
)

func insert(db *bolt.DB){
    _ = db.Update(func(tx *bolt.Tx) error {
        bucket := tx.Bucket([]byte("MyBucket"))
        c := bucket.Cursor()
        k, _ := c.Last()
        i, _ := strconv.Atoi(string(k[:]))
        str := fmt.Sprintf("%d", i + 1)
        bucket.Put([]byte(str), []byte("Test"))
        return nil
    })
}

func main() {
    db, err := bolt.Open("my.db", 0600, nil)
    if err != nil {
        log.Fatal(err)
    }
    
    insert(db)
    db.Update(func(tx *bolt.Tx) error {
        tx.CreateBucketIfNotExists([]byte("MyBucket"))
        return nil
    });
    
    db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte("MyBucket"))
        b.ForEach(func(k, v []byte) error {
            fmt.Printf("key=%s, value=%s\n", k, v)
            return nil
        })
        return nil
    })
    
    defer db.Close()
}