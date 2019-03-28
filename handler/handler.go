package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/lib/pq"

	redigo "github.com/garyburd/redigo/redis"
	"github.com/julienschmidt/httprouter"
	"github.com/mnirfan/bigproject/storage"
)

type Person struct {
	ID          int         `json:"id"`
	Name        *string     `json:"name"`
	MSISDN      int64       `json:"msisdn"`
	Email       string      `json:"email"`
	BirthDate   *string     `json:"birthDate"`
	CreatedTime time.Time   `json:"createdTime"`
	UpdatedTime pq.NullTime `json:"updatedTime"`
}

type Respon struct {
	Counter int      `json:"counter"`
	Users   []Person `json:"user"`
}

func GetIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	q := r.URL.Query()
	name := q.Get("search")
	db, err := storage.ConnectSQL()
	if err != nil {
		log.Println("Err: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	redisPool := redigo.Pool{
		IdleTimeout: 10,
		MaxActive:   30,
		MaxIdle:     240,
		Dial: func() (redigo.Conn, error) {
			return redigo.Dial("tcp", "localhost:6379")
		},
	}

	redisConn := redisPool.Get()

	// redisConn := storage.ConnectRedis()
	defer redisConn.Close()

	keyCache := "class:guess:counter"

	_, err = redisConn.Do("INCR", keyCache)
	if err != nil {
		log.Println("Err: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resultCache, err := redisConn.Do("GET", keyCache)
	if err != nil {
		log.Println("Err: 1", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	counter, err := redigo.Int(resultCache, err)

	rows, err := db.Query(fmt.Sprintf(`
		SELECT
			user_id,
			user_name,
			msisdn,
			user_email,
			birth_date,
			create_time,
			update_time
		FROM
			ws_user
		WHERE
			user_name
		ILIKE
			'%%%v%%'
		LIMIT
			%v
	`, name, 10))

	if err != nil {
		log.Println("Err: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var result []Person
	for rows.Next() {
		temp := Person{}
		err = rows.Scan(
			&temp.ID,
			&temp.Name,
			&temp.MSISDN,
			&temp.Email,
			&temp.BirthDate,
			&temp.CreatedTime,
			&temp.UpdatedTime,
		)
		if err != nil {
			log.Println("Err: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		result = append(result, temp)
	}
	resp := Respon{
		Counter: counter,
		Users:   result,
	}
	jsonData, err := json.Marshal(resp)
	if err != nil {
		log.Println("Err: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// fmt.Println(w.Header())
	w.WriteHeader(200)
	w.Write(jsonData)
}

func OptionsIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// fmt.Println(w.Header())
	w.WriteHeader(200)
}
