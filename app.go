package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/mnirfan/bigproject/handler"
)

func main() {
	router := httprouter.New()
	router.OPTIONS("/", handler.OptionsIndex)
	router.GET("/", handler.GetIndex)

	log.Fatal(http.ListenAndServe(":8080", router))
}
