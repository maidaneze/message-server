package controllers

import "net/http"

func (h Handler) Setup() *http.Server {
	srv := &http.Server{Addr: ":8080"}
	http.HandleFunc("/check", h.Check)
	http.HandleFunc("/users", h.CreateUser)
	http.HandleFunc("/login", h.LoginUser)
	http.HandleFunc("/messages", h.MessagesHandler)
	return srv
}
