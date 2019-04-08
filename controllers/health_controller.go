package controllers

import (
	"encoding/json"
	"net/http"
)

//Returns 200 if the database check succeeds and 500 otherwise

func (h Handler) Check(w http.ResponseWriter, r *http.Request) {
	if err := h.Db.CheckConnection(); err != nil {
		http.Error(w, "DB connection error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]string{"health": "ok"}); err != nil {
		http.Error(w, "Write error", http.StatusInternalServerError)
	}
}
