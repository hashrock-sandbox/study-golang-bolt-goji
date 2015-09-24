package main

import (
    "log"
    "fmt"
    "github.com/boltdb/bolt"
    "strconv"
    "net/http"
	"github.com/goji/param"    
    "github.com/zenazn/goji"
    "github.com/zenazn/goji/web"
    "encoding/json"
)


func insert(db *bolt.DB, contents string){
    _ = db.Update(func(tx *bolt.Tx) error {
        bucket := tx.Bucket([]byte("MyBucket"))
        c := bucket.Cursor()
        k, _ := c.Last()
        i, _ := strconv.Atoi(string(k[:]))
        str := fmt.Sprintf("%06d", i + 1)
        bucket.Put([]byte(str), []byte(contents))
        return nil
    })
}

type Post struct{
    PostId string
    Contents string
}

func insertPost(w http.ResponseWriter, r *http.Request){
    db, err := bolt.Open("my.db", 0600, nil)
    if err != nil {
        log.Fatal(err)
    }    

    var p Post
    r.ParseForm()
	err = param.Parse(r.Form, &p)
    if err != nil {
        log.Fatal(err)
    }
    insert(db, p.Contents)
    
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, "{\"status\": \"success\"}")
    defer db.Close()
}



func readPosts(c web.C, w http.ResponseWriter, r *http.Request) {
    db, err := bolt.Open("my.db", 0600, nil)
    if err != nil {
        log.Fatal(err)
    }
    
    var posts = make([]Post, 0, 0)
    
    db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte("MyBucket"))
        b.ForEach(func(k, v []byte) error {
            var p = Post{PostId: string(k[:]), Contents: string(v[:])}
            posts = append(posts, p)
            return nil
        })
        return nil
    })
    defer db.Close()
    
    w.Header().Set("Content-Type", "application/json")
    encoder := json.NewEncoder(w)
    encoder.Encode(posts)
}

func createDB(){
    db, err := bolt.Open("my.db", 0600, nil)
    if err != nil {
        log.Fatal(err)
    }
    db.Update(func(tx *bolt.Tx) error {
        tx.CreateBucketIfNotExists([]byte("MyBucket"))
        return nil
    });
    defer db.Close()
}

func main() {
    createDB()
    goji.Get("/api/posts/", readPosts)
    goji.Post("/api/posts/", insertPost)
    goji.Get("/*", http.FileServer(http.Dir("./static/")))
    goji.Serve()
}