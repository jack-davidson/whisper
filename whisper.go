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
	http.HandleFunc("/me", Me)
	http.ListenAndServe(":3000", nil)
}
