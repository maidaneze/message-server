package main

import (
	"fmt"
	"github.com/maidaneze/message-server/controllers"
	"github.com/maidaneze/message-server/dao"
	"log"
)

func main() {
	db, err := dao.OpenSqlite3Database("challenge.db")
	if err != nil {
		log.Fatal("Unable to open DB")
	}


	h := controllers.Handler{Db: db}
	server := h.Setup()
	fmt.Println("Server started!!")
	log.Fatal(server.ListenAndServe())
}
