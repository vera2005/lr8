package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "catjkm8800"
	dbname   = "query"
)

type Handlers struct {
	dbProvider DatabaseProvider
}

type DatabaseProvider struct {
	db *sql.DB
}

// Обработчики HTTP-запросов
func (h *Handlers) GetQuery(w http.ResponseWriter, r *http.Request) {
	msg, err := h.dbProvider.SelectQuery()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello " + msg + "!"))
}

func (h *Handlers) PostQuery(w http.ResponseWriter, r *http.Request) {
	nameInput := r.URL.Query().Get("name") // ради разнообразия поработаем с Query-параметром, как в изначальном условии лр6
	if nameInput == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing 'name' query parameter"))
	}
	err := h.dbProvider.InsertQuery(nameInput)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handlers) PutQuery(w http.ResponseWriter, r *http.Request) {
	nameInput := r.URL.Query().Get("name") // ради разнообразия поработаем с Query-параметром, как в изначальном условии лр6
	if nameInput == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing 'name' query parameter"))
		return
	}
	err := h.dbProvider.UpdateQuery(nameInput)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// Методы для работы с базой данных
func (dbp *DatabaseProvider) SelectQuery() (string, error) {
	var msg string
	row := dbp.db.QueryRow("SELECT name FROM query ORDER BY id DESC LIMIT 1")
	err := row.Scan(&msg)
	if err != nil {
		return "", err
	}
	return msg, nil
}

func (dbp *DatabaseProvider) UpdateQuery(n string) error {
	_, err := dbp.db.Exec("UPDATE query SET name = $1 WHERE id = (SELECT MAX(id) FROM query)", n)
	if err != nil {
		return err
	}
	return nil
}

func (dbp *DatabaseProvider) InsertQuery(n string) error {
	_, err := dbp.db.Exec("INSERT INTO query (name) VALUES ($1)", n)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	address := flag.String("address", "127.0.0.1:8081", "адрес для запуска сервера")
	flag.Parse()
	// Формирование строки подключения для postgres
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// Создание соединения с сервером postgres
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		fmt.Println("NO! 2")
	}
	fmt.Println("Connected!")
	// Создаем провайдер для БД с набором методов
	dp := DatabaseProvider{db: db}
	// Создаем экземпляр структуры с набором обработчиков
	h := Handlers{dbProvider: dp}

	// Регистрируем обработчики
	http.HandleFunc("/get", h.GetQuery)
	http.HandleFunc("/post", h.PostQuery)
	http.HandleFunc("/put", h.PutQuery)

	// Запускаем веб-сервер на указанном адресе
	err = http.ListenAndServe(*address, nil)
	if err != nil {
		log.Fatal(err)
	}

}
