package main

import (
	"database/sql"
	"encoding/json"
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
	dbname   = "count"
)

type Handlers struct {
	dbProvider DatabaseProvider
}
type DatabaseProvider struct {
	db *sql.DB
}

// Обработчики HTTP-запросов
func (h *Handlers) GetCount(w http.ResponseWriter, r *http.Request) {
	msg, err := h.dbProvider.SelectCount()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(msg))
}

func (h *Handlers) PostCount(w http.ResponseWriter, r *http.Request) {
	input := struct {
		Val float32 `json:"val"`
	}{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	err = h.dbProvider.InsertCount(input.Val)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handlers) PutCount(w http.ResponseWriter, r *http.Request) {

	input := struct {
		Val float32 `json:"val"`
	}{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	err = h.dbProvider.UpdateCount(input.Val)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.WriteHeader(http.StatusCreated)
}

// Методы для работы с базой данных
func (dbp *DatabaseProvider) SelectCount() (string, error) {
	var msg string
	row := dbp.db.QueryRow("SELECT summa FROM count ORDER BY id DESC LIMIT 1")
	err := row.Scan(&msg)
	if err != nil {
		return "", err
	}
	return msg, nil
}

func (dbp *DatabaseProvider) InsertCount(v float32) error {
	_, err := dbp.db.Exec("INSERT INTO count (val, summa) VALUES ($1, $1+(SELECT summa FROM count ORDER BY id DESC LIMIT 1))", v)
	if err != nil {
		return err
	}

	return nil
}

func (dbp *DatabaseProvider) UpdateCount(v float32) error {
	_, err := dbp.db.Exec("UPDATE count SET val = $1, summa = (val + (SELECT summa FROM count WHERE id = ((SELECT MAX(id) FROM count) - 1))) WHERE id = (SELECT MAX(id) FROM count)", v)
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
	http.HandleFunc("/get", h.GetCount)
	http.HandleFunc("/post", h.PostCount)
	http.HandleFunc("/put", h.PutCount)

	// Запускаем веб-сервер на указанном адресе
	err = http.ListenAndServe(*address, nil)
	if err != nil {
		log.Fatal(err)
	}

}
