package controllers

import (
	"bytes"
	"github.com/maidaneze/message-server/model"
	"github.com/maidaneze/message-server/services/auth"
	"github.com/maidaneze/message-server/services/messages"
	"encoding/json"
	"net/http"
)

//Wrapper for the insertMessage handler and getMessage handler
//Executes the getMessage handler if the resource is GET
//Executes the insertMessage handler if the resource is POST
//Returns 404 otherwise
func (h Handler) MessagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.getMessage(w, r)
	} else if r.Method == "POST" {
		h.insertMessage(w, r)
	} else {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
}

//Sends a new message from the sender to the recipient
//Returns 400 if the request dto is invalid
//Returns 401 if the token header doesn't correspond to the sender or is invalid
//Returns 500 if the message couldn't be sent
//Returns 200 otherwise
func (h Handler) insertMessage(w http.ResponseWriter, r *http.Request) {

	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	defer r.Body.Close()

	var postMessageRequestDTO model.PostMessageRequestDTO

	//Check body
	var err error
	if err := json.Unmarshal(buf.Bytes(), &postMessageRequestDTO); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	var dto model.MessageDTO
	if dto, err = messages.UnmarshallMessageContent(postMessageRequestDTO); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	//Validate dto
	if !messages.ValidMessageDto(dto) {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	//Validate authorization Header
	token, valid := auth.ValidateTokenHeader(r)
	if !valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	//Validate recipient Exists
	found, err := h.Db.CheckUserExists(dto.RecipientId)
	if !found && err == nil {
		http.Error(w, "Recipient doen't exist", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "Error sending message", http.StatusInternalServerError)
		return
	}

	//Validate sender Exists
	found, err = h.Db.CheckUserExists(dto.SenderId)
	if !found && err == nil {
		http.Error(w, "Sender doen't exist", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "Error sending message", http.StatusInternalServerError)
		return
	}

	//Get user tokens
	tokens, err := h.Db.GetTokens(dto.SenderId)

	if err != nil {
		http.Error(w, "Error sending message", http.StatusInternalServerError)
		return
	}

	//Validate user token
	valid = auth.ValidateAuthorizedUser(token, tokens)
	if !valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	//Insert into messages
	dto, err = h.Db.InsertMessage(dto)
	if err != nil {
		http.Error(w, "Error sending message", http.StatusInternalServerError)
		return
	}

	//Marshall response
	if err := json.NewEncoder(w).Encode(messages.GetPostMessageResponseDTO(dto)); err != nil {
		http.Error(w, "Error sending message", http.StatusInternalServerError)
	}
}

//Gets all messages for the recipient starting from the given messageid and up to limit number of messages (Default 100)
//Returns 400 if the request queryparams are invalid
//Returns 401 if the token header doesn't correspond to the sender or is invalid
//Returns 500 if the message couldn't be recovered
//Returns 200 otherwise
func (h Handler) getMessage(w http.ResponseWriter, r *http.Request) {

	//Get Query params
	recipientid, messageid, limit, err := messages.ParseGetMessageQueryParams(r)

	if err != nil {
		http.Error(w, "Invalid Params", http.StatusBadRequest)
		return
	}

	//Validate query params
	if !(recipientid > 0 && messageid > 0 && limit > 0) {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	//Validate authorization Header
	token, valid := auth.ValidateTokenHeader(r)
	if !valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	//Validate recipient Exists
	found, err := h.Db.CheckUserExists(recipientid)
	if !found && err == nil {
		http.Error(w, "Recipient doen't exist", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "Error getting messages", http.StatusInternalServerError)
		return
	}

	//Get user tokens
	tokens, err := h.Db.GetTokens(recipientid)

	if err != nil {
		http.Error(w, "Error getting messages", http.StatusInternalServerError)
		return
	}

	//Validate user token
	valid = auth.ValidateAuthorizedUser(token, tokens)
	if !valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	//Get messages
	getMessages, err := h.Db.GetMessages(recipientid, messageid, limit)
	if err != nil {
		http.Error(w, "Error getting messages", http.StatusInternalServerError)
		return
	}

	messagesResponse := messages.ParseMessages(getMessages)

	//Marshall response
	if err := json.NewEncoder(w).Encode(map[string][]model.MessageResponse{"messages": messagesResponse}); err != nil {
		http.Error(w, "Error getting messages", http.StatusInternalServerError)
	}
}
