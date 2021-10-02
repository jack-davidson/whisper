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

func lookupUserID(name string) *string {
	db, err := sql.Open("postgres", DBDataSourceName)
	if err != nil {
		fmt.Println("Failed to connect to db")
	}
	defer db.Close()
	q, err := db.Query("select id from users where name=$1;", name)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	if q.Next() {
		var id string
		q.Scan(&id)
		return &id
	}
	return nil
}

func authenticate(name string, auth string, then func(db *sql.DB, id string)) {
	db, err := sql.Open("postgres", DBDataSourceName)
	if err != nil {
		fmt.Println("Failed to connect to db")
	}
	defer db.Close()
	q, err := db.Query("select id from users where name=$1 and passwordhash=$2;", name, hash(auth))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if q.Next() {
		var id string
		q.Scan(&id)
		then(db, id)
	}
}

func NewUser(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", DBDataSourceName)
	if err != nil {
		fmt.Println("Failed to connect to db")
	}
	defer db.Close()

	name := r.Header.Get("Name")
	password := r.Header.Get("Password")
	q, err := db.Query("select * from users where name=$1;", name)
	if !q.Next() {
		auth := genHash()
		db.Exec("insert into users (name, passwordhash) values ($1, $2);", name, hash(password))
		fmt.Fprintln(w, auth)
	} else {
		w.WriteHeader(http.StatusConflict)
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	authenticate(r.Header.Get("Name"), r.Header.Get("Password"), func(db *sql.DB, id string) {
		db.Exec("delete from users where id=$1;", id)
	})
}

func Messages(w http.ResponseWriter, r *http.Request) {
	authenticate(r.Header.Get("Name"), r.Header.Get("Password"), func(db *sql.DB, id string) {
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
	})
}

func Message(w http.ResponseWriter, r *http.Request) {
	authenticate(r.Header.Get("Name"), r.Header.Get("Password"), func(db *sql.DB, id string) {
		text := r.Header.Get("Text")
		forUserID := r.Header.Get("For")
		db.Exec("insert into messages (text, foruser, fromuser) values ($1, $2, $3);", text, forUserID, id)
	})
}

func User(w http.ResponseWriter, r *http.Request) {
	id := lookupUserID(r.Header.Get("Name"))
	if id != nil {
		fmt.Fprintln(w, *id)
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})
	http.HandleFunc("/newuser", NewUser)
	http.HandleFunc("/deleteuser", DeleteUser)
	http.HandleFunc("/messages", Messages)
	http.HandleFunc("/message", Message)
	http.HandleFunc("/user", User)
	http.ListenAndServe(":3000", nil)
}
