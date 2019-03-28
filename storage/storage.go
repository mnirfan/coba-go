package storage

import (
	"database/sql"
	"fmt"

	redigo "github.com/garyburd/redigo/redis"
)

func ConnectSQL() (*sql.DB, error) {
	// db, err := sql.Open("postgres", fmt.Sprintf(`
	//     dbname=bigproject
	//     user=opbdev
	//     password=opbdev
	//     host=localhost
	//     sslmode=disable
	// `))
	db, err := sql.Open("postgres", fmt.Sprintf(`
	    dbname=tokopedia-user
	    user=tkpdtraining
	    password=trainingyangbeneryah
	    host=devel-postgre.tkpd
	    sslmode=disable
	`))
	return db, err
}

func ConnectRedis() redigo.Conn {
	redisPool := redigo.Pool{
		IdleTimeout: 10,
		MaxActive:   30,
		MaxIdle:     240,
		Dial: func() (redigo.Conn, error) {
			return redigo.Dial("tcp", "localhost:6379")
		},
	}

	redisConn := redisPool.Get()
	return redisConn
}
