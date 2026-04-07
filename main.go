package main

import (
	"database/sql"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

//go:embed index.html
var indexHTML embed.FS

var db *sql.DB
var tmpl *template.Template

func initDB() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	if port == "" {
		port = "5432"
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("DB 연결 실패: ", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("DB Ping 실패: ", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS visitors (
			id SERIAL PRIMARY KEY,
			ip VARCHAR(50),
			visited_at TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		log.Fatal("테이블 생성 실패: ", err)
	}

	log.Println("DB 연결 성공!")
}

func index(w http.ResponseWriter, r *http.Request) {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	} else {
		// X-Forwarded-For can contain multiple IPs: "client, proxy1, proxy2"
		// The first one is the real client IP
		ip = strings.TrimSpace(strings.Split(ip, ",")[0])
	}

	_, err := db.Exec("INSERT INTO visitors (ip) VALUES ($1)", ip)
	if err != nil {
		log.Println("방문자 저장 실패: ", err)
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM visitors").Scan(&count)
	if err != nil {
		log.Println("방문자 수 조회 실패: ", err)
	}

	data := struct {
		Count int
		IP    string
	}{Count: count, IP: ip}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, data)
}

func main() {
	var err error
	tmpl, err = template.ParseFS(indexHTML, "index.html")
	if err != nil {
		log.Fatal("템플릿 로드 실패: ", err)
	}

	initDB()
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	http.HandleFunc("/", index)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
