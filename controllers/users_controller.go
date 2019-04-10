package controllers

import (
	"bytes"
	"github.com/maidaneze/message-server/model"
	"github.com/maidaneze/message-server/services/auth"
	"github.com/maidaneze/message-server/services/passwords"
	"github.com/maidaneze/message-server/services/users"
	"encoding/json"
	"net/http"
)

//Creates a new user for the given username and password
//Returns 400 if the request dto is invalid
//Returns 409 if the user already exists
//Returns 500 if the user couldn't be created
//Returns 200 otherwise

func (h Handler) CreateUser(w http.ResponseWriter, r *http.Request) {

	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	defer r.Body.Close()

	var usersRequestDTO model.UserRequestDTO

	//Unmarshall body

	var err error
	if err = json.Unmarshal(buf.Bytes(), &usersRequestDTO); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	//Validate dto

	if !users.ValidateUsersRequestDTO(usersRequestDTO) {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	//Check If user exists

	var found bool
	if _, found, err = h.Db.GetUser(usersRequestDTO.Username); found && err == nil {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	if !found && err != nil {
		http.Error(w, "Error generating user", http.StatusInternalServerError)
		return
	}

	//Create User

	var user model.User
	if user, err = users.CreateUser(usersRequestDTO.Username, usersRequestDTO.Password); err != nil {
		http.Error(w, "Error generating user", http.StatusInternalServerError)
		return
	}

	//Post user

	user, err = h.Db.InsertUser(user)
	if err != nil {
		http.Error(w, "Error generating user", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(model.UserResponseDTO{user.Userid}); err != nil {
		http.Error(w, "Write error", http.StatusInternalServerError)
	}
}

//Generates a new access token for the given username and password
//Returns 400 if the request dto is invalid
//Returns 401 if the user doesn't exist or the username/password combination is invalid
//Returns 500 if the token couldn't be insertd
//Returns 200 otherwise the userid and the token

func (h Handler) LoginUser(w http.ResponseWriter, r *http.Request) {

	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	defer r.Body.Close()

	var usersRequestDTO model.UserRequestDTO

	//Check body

	if err := json.Unmarshal(buf.Bytes(), &usersRequestDTO); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	//Validate dto

	if !users.ValidateUsersRequestDTO(usersRequestDTO) {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	//Get user

	user, found, err := h.Db.GetUser(usersRequestDTO.Username)

	if err == nil && !found {
		http.Error(w, "Invalid username password combination", http.StatusUnauthorized)
		return
	}

	if err != nil && !found {
		http.Error(w, "Error logging in", http.StatusInternalServerError)
		return
	}

	//Validate user

	if !passwords.ValidUserPassword(user, usersRequestDTO.Password) {
		http.Error(w, "Invalid username password combination", http.StatusUnauthorized)
		return
	}

	//Generate token

	token, err := auth.GenerateToken()

	if err != nil {
		http.Error(w, "Error logging in", http.StatusInternalServerError)
		return
	}

	err = h.Db.InsertToken(user.Userid, token)

	if err != nil {
		http.Error(w, "Error logging in", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(model.LoginResponseDTO{user.Userid, token.Uuid}); err != nil {
		http.Error(w, "Error logging in", http.StatusInternalServerError)
	}
}
