package main

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)

type from_db_url struct {
	url string
}

var database *sql.DB

func AddtoDB(url string, code string) {
	if GetFromDB(code) == "sql: no rows in result set" {
		_, err := database.Exec("INSERT INTO shortener (url, code) VALUES ($1, $2)", url, code)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("Error")
	}

}

func GetFromDB(code string) string {

	row := database.QueryRow("SELECT url FROM shortener where code = $1", code)

	answer := from_db_url{}
	err := row.Scan(&answer.url)

	if err != nil {
		return err.Error()
	}
	return answer.url

}

func MakeShort(w http.ResponseWriter, r *http.Request) {

	url := r.URL.Query().Get("url")
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(url)))[:10]

	AddtoDB(url, string(hash))
	w.Write([]byte("Kodirovano v: "))
	w.Write([]byte(hash + "\n"))
	w.Write([]byte("Декодировать можно по адресу: http://127.0.0.1:8000/s/" + hash))
}

func DecodeURL(w http.ResponseWriter, r *http.Request) {
	hash_code := r.RequestURI[3:len(r.RequestURI)]

	fmt.Printf("hash is: %s\n", hash_code)
	w.Write([]byte("Dekodirovano v: "))
	w.Write([]byte(GetFromDB(hash_code)))

}

func main() {
	connStr := "user=postgres password=Rkhg2s9121096 dbname=shortener sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	database = db
	defer db.Close()
	db.SetMaxOpenConns(100)

	http.HandleFunc("/s/", DecodeURL)
	http.HandleFunc("/a/", MakeShort)
	http.ListenAndServe(":8000", nil)

}
