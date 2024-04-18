package main

import (
	"crypto/sha256"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
)

var redisdb *redis.Client

func AddtoDB() {

	opt, err := redis.ParseURL("redis://localhost:6379/0")
	if err != nil {
		panic(err)
	}

	fmt.Println("addr is", opt.Addr)
	fmt.Println("db is", opt.DB)
	fmt.Println("password is", opt.Password)
	//client := redis.NewClient(opt)

	//_ := redis.NewClient(opt)

}
func MakeShort(w http.ResponseWriter, r *http.Request) {
	ss := "https://gobyexample.com"

	resp, err := http.Get(ss)

	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		fmt.Println(resp.StatusCode)
	} else {
		fmt.Printf("Status code of website: %d\n", resp.StatusCode)
	}

	url := r.URL.Query().Get("url")
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(url)))[:10]

	fmt.Printf("ssilka: %s\n", url)
	fmt.Printf("hashed to: %s\n", hash)
}

func main() {
	AddtoDB()
	http.HandleFunc("/a/", MakeShort)
	http.ListenAndServe(":8000", nil)

}
