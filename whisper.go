package main

import (
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func readDBPassword() string {
	b, err := os.ReadFile("dbpasswd")
	if err != nil {
		fmt.Println("Failed to read dbpasswd")
	}
	return string(b)
}

const (
	DBHost = "localhost"
	DBPort = 5432
	DBUser = "postgres"
	DB     = "whisper"
)

var (
	DBPasswd         = readDBPassword()
	DBDataSourceName = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", DBHost, DBPort, DBUser, DBPasswd, DB)
)

func genHash() string {
	auth := make([]byte, 64)
	rand.Seed(time.Now().UnixNano())
	rand.Read(auth)
	return hex.EncodeToString(auth)
}

func hash(s string) string {
	authHash := sha512.Sum512([]byte(s))
	return hex.EncodeToString(authHash[:])
}

func NewUser(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", DBDataSourceName)
	if err != nil {
		fmt.Println("Failed to connect to db")
	}
	defer db.Close()

	name := r.Header.Get("Name")
	auth := genHash()
	fmt.Fprintln(w, auth)
	db.Exec("insert into users (name, authhash) values ($1, $2);", name, hash(auth))
}

func Messages(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", DBDataSourceName)
	if err != nil {
		fmt.Println("Failed to connect to db")
	}
	defer db.Close()

	name := r.Header.Get("Name")
	auth := r.Header.Get("Auth")

	q, err := db.Query("select id from users where name=$1 and authhash=$2;", name, hash(auth))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if q.Next() {
		var id string
		q.Scan(&id)
		q, err := db.Query("select * from messages where foruser=$1;", id)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		for q.Next() {
			var id, text, foruser, fromuser interface{}
			q.Scan(&id, &foruser, &fromuser, &text)
			fmt.Fprintf(w, "%s '%s' %s %s\n", id, text, foruser, fromuser)
		}
	}
}

func Message(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", DBDataSourceName)
	if err != nil {
		fmt.Println("Failed to connect to db")
	}
	defer db.Close()

	name := r.Header.Get("Name")
	auth := r.Header.Get("Auth")
	text := r.Header.Get("Text")

	q, err := db.Query("select id from users where name=$1 and authhash=$2;", name, hash(auth))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if q.Next() {
		var fromuser string
		q.Scan(&fromuser)
		db.Exec("insert into messages (text, foruser, fromuser) values ($1, $2, $3);", text, fromuser, fromuser)
	}
}

func Me(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", DBDataSourceName)
	if err != nil {
		fmt.Println("Failed to connect to db")
	}
	defer db.Close()

	name := r.Header.Get("Name")
	auth := r.Header.Get("Auth")

	q, err := db.Query("select * from users where name=$1 and authhash=$2;", name, hash(auth))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if q.Next() {
		fmt.Fprintf(w, "success\n")
	} else {
		fmt.Fprintf(w, "failed\n")
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})
	http.HandleFunc("/newuser", NewUser)
	http.HandleFunc("/messages", Messages)
	http.HandleFunc("/message", Message)
	http.HandleFunc("/me", Me)
	http.ListenAndServe(":3000", nil)
}
