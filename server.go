package main

import (
	"challenge/controllers"
	"challenge/dao"
	"log"
)

func main() {
	db, err := dao.OpenSqlite3Database("challenge.db")
	if err != nil {
		log.Fatal("Unable to open DB")
	}
	h := controllers.Handler{Db: db}
	server := h.Setup()

	log.Fatal(server.ListenAndServe())
}
