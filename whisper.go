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

func genAuth() string {
	auth := make([]byte, 64)
	rand.Seed(time.Now().UnixNano())
	rand.Read(auth)
	return hex.EncodeToString(auth)
}

func NewUser(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", DBDataSourceName)
	if err != nil {
		fmt.Println("Failed to connect to db")
	}
	defer db.Close()

	name := r.Header.Get("Name")
	auth := genAuth()
	fmt.Fprintln(w, auth)
	authHash := sha512.Sum512([]byte(auth))
	db.Exec(fmt.Sprintf("insert into users (name, authhash) values ('%s', '%s');", name, hex.EncodeToString(authHash[:])))
}

func Me(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", DBDataSourceName)
	if err != nil {
		fmt.Println("Failed to connect to db")
	}
	defer db.Close()

	name := r.Header.Get("Name")
	auth := r.Header.Get("Auth")
	authHash := sha512.Sum512([]byte(auth))
	q, err := db.Query(fmt.Sprintf("select * from users where name='%s' and authhash='%s';", name, hex.EncodeToString(authHash[:])))
	if err != nil {
		fmt.Println("Query failed")
	}
	if q.Next() {
		fmt.Fprintf(w, "Auth Success\n")
	} else {
		fmt.Fprintf(w, "Failed to Authenticate as: %s, with Auth: %s\n", name, auth)
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})
	http.HandleFunc("/newuser", NewUser)
	http.HandleFunc("/me", Me)
	http.ListenAndServe(":3000", nil)
}
