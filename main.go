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
	res, err := database.Exec("INSERT INTO shortener (url, code) VALUES ($1, $2)", url, code)

	if err != nil {
		panic(err)
	}
	fmt.Print("Добавлено")
	fmt.Println(res.RowsAffected())

}

func GetFromDB(code string) {

	row := database.QueryRow("SELECT url FROM shortener where code = $1", code)

	answer := from_db_url{}
	err := row.Scan(&answer.url)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(answer.url)
	// kol-vo konnektov - 11
	/*

		var count int
		errik := database.QueryRow("SELECT COUNT(*) FROM pg_stat_activity").Scan(&count)

		if errik != nil {
			fmt.Println(errik)
		}
		fmt.Printf("There are %d connections in the database\n", count)
	*/

}

func MakeShort(w http.ResponseWriter, r *http.Request) {

	url := r.URL.Query().Get("url")
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(url)))[:10]

	AddtoDB(url, string(hash))
	fmt.Println("Dobavleno: ")
	fmt.Printf("ssilka: %s\n", url)
	fmt.Printf("hashed to: %s\n", hash)
}

func DecodeURL(w http.ResponseWriter, r *http.Request) {
	hash_code := r.RequestURI[3:len(r.RequestURI)]

	fmt.Printf("hash is: %s\n", hash_code)
	fmt.Println("Url is: ")
	GetFromDB(hash_code)

}

func main() {
	connStr := "user=postgres password=Rkhg2s9121096 dbname=shortener sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	database = db
	defer db.Close()
	db.SetMaxOpenConns(20)

	http.HandleFunc("/s/", DecodeURL)
	http.HandleFunc("/a/", MakeShort)
	http.ListenAndServe(":8000", nil)

}
