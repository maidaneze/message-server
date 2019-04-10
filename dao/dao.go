package dao

import "github.com/maidaneze/message-server/model"

type DB interface {

	//Verifies the connection to the database
	//Returns error in case of failiure and nil in case of success

	CheckConnection() error

	//Inserts a new user into the users tables
	//Returns error in case of failiure and the inserted user in case of success

	InsertUser(user model.User) (model.User, error)

	//Gets the given user data
	//Returns true if an user was found, false otherwise
	//Returns error in case of failiure and the user in case of success

	GetUser(username string) (model.User, bool, error)

	//Checks if the corresponding user id is registered
	//Returns true if an user was found, false otherwise
	//Returns error in case of failiure

	CheckUserExists(userid int64) (bool, error)

	//Inserts the given token into the tokens table for the requested userid
	//Returns error in case of failiure and nil in case of success

	InsertToken(userid int64, token model.Token) error

	//Recovers the access tokens for the requested userid (at most 2)
	//Returns error in case of failiure and nil in case of success

	GetTokens(userid int64) ([]model.Token, error)

	//Inserts the given message into the messages table
	//Returns error in case of failiure and nil in case of success

	InsertMessage(message model.MessageDTO) (model.MessageDTO, error)

	//Recovers the messages for the requested recipient id that have a messageid greater than the given messageid
	//Returns at most the given limit of messages
	//Returns error in case of failiure and nil in case of success

	GetMessages(recipientId int64, messageId int64, limit int64) ([]model.MessageDTO, error)
}
