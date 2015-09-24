package main

import (
    "log"
    "fmt"
    "github.com/boltdb/bolt"
    "strconv"
)

func insert(db *bolt.DB){
    _ = db.Update(func(tx *bolt.Tx) error {
        bucket, _ := tx.CreateBucketIfNotExists([]byte("MyBucket"))
        bucket.Put([]byte("1"), []byte("Test"))
        return nil
    })
}

/*
func newPost(db *bolt.DB){
    db.Update(func(tx *bolt.Tx) error {
        bucket, err := tx.CreateBucketIfNotExists([]byte("MyBucket"))
        if err != nil {
            log.Fatal(err)
        }    
        c := bucket.Cursor()
        k, _ := c.Last()
        
        fmt.Printf("last=%s\n", k)
        fmt.Printf("last=%s\n", len(k))

        num := binary.BigEndian.Uint64(k)
        fmt.Printf("last=%s\n", num)

        //fmt.Printf("last=%s\n", k + 1)
        //data := binary.BigEndian.Uint64(k)
        //fmt.Println(data)
        err = bucket.Put([]byte(num), []byte("value1"))
        
        if err != nil {
            log.Fatal(err)
        }

        return nil
    })    
}
*/

func main() {
    // Open the my.db data file in your current directory.
    // It will be created if it doesn't exist.
    db, err := bolt.Open("my.db", 0600, nil)
    if err != nil {
        log.Fatal(err)
    }
    
    //newPost(db)
    insert(db)
    
    
    db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte("MyBucket"))
        b.ForEach(func(k, v []byte) error {
            fmt.Printf("key=%s, value=%s\n", k, v)
            return nil
        })
        
        c := b.Cursor()
        k, _ := c.Last()
        
        i, _ := strconv.Atoi(string(k[:]))
        fmt.Printf("next=%s\n", fmt.Sprint(i + 1))
        return nil
    })
    
    defer db.Close()
}