package dao

import (
	"github.com/maidaneze/message-server/utils"
	"database/sql"
	"errors"
	"time"

	"github.com/maidaneze/message-server/model"
)

type SqliteDB struct {
	db *sql.DB
}

func OpenSqlite3Database(dbname string) (SqliteDB, error) {
	var db *sql.DB
	var err error
	err = utils.Retry(func() error {
		db, err = sql.Open("sqlite3", dbname)
		return err
	}, 2, time.Millisecond*20)
	return SqliteDB{db}, err
}

//Wrapper function for "checkConnection"
//Executes insertUser with a retry

func (sqlite SqliteDB) CheckConnection() error {
	var err error
	err = utils.Retry(func() error {
		err = sqlite.checkConnection()
		return err
	}, 2, time.Millisecond*20)
	return err
}

//Verifies the connection to the database
//Returns error in case of failiure and nil in case of success

func (sqlite SqliteDB) checkConnection() error {
	var res int
	if err := sqlite.db.QueryRow("SELECT 1").Scan(&res); err != nil {
		return err
	}

	if res != 1 {
		return errors.New("Unexpected query result")
	}
	return nil
}

//Wrapper function for insertUser
//Executes insertUser with a retry

func (sqlite SqliteDB) InsertUser(user model.User) (model.User, error) {
	var insertedUser model.User
	var err error
	err = utils.Retry(func() error {
		insertedUser, err = sqlite.insertUser(user)
		return err
	}, 2, time.Millisecond*20)
	return insertedUser, err
}

//Inserts a new user into the users tables
//Returns error in case of failiure and the inserted user in case of success

func (sqlite SqliteDB) insertUser(user model.User) (model.User, error) {
	//Insert user
	var insertResult sql.Result
	var err error
	if insertResult, err = sqlite.db.Exec(insertNewUserQuey, user.Username, user.Password, user.Salt); err != nil {
		return model.User{}, err
	}

	id, err := insertResult.LastInsertId()
	user.Userid = id
	return user, err
}

//Wrapper function for getUser
//Executes getUser with a retry

func (sqlite SqliteDB) GetUser(username string) (model.User, bool, error) {
	var getUser model.User
	var err error
	var found bool
	err = utils.Retry(func() error {
		getUser, found, err = sqlite.getUser(username)
		return err
	}, 2, time.Millisecond*20)
	return getUser, found, err
}

//Gets the given user data
//Returns true if an user was found, false otherwise
//Returns error in case of failiure and the user in case of success

func (sqlite SqliteDB) getUser(username string) (model.User, bool, error) {
	var userid int64
	var password string
	var password_salt string

	if err := sqlite.db.QueryRow(getFromUsersQuery, username).Scan(&userid, &password, &password_salt); err != nil {
		if err == sql.ErrNoRows {
			return model.User{}, false, nil
		} else {
			return model.User{}, false, err
		}
	}

	user := model.User{userid, username, password, password_salt}
	return user, true, nil
}

//Wrapper function for checkUserExists
//Executes checkUserExists with a retry

func (sqlite SqliteDB) CheckUserExists(userid int64) (bool, error) {
	var err error
	var found bool
	err = utils.Retry(func() error {
		found, err = sqlite.checkUserExists(userid)
		return err
	}, 2, time.Millisecond*20)
	return found, err
}

//Checks if the corresponding user id is registered
//Returns true if an user was found, false otherwise
//Returns error in case of failure

func (sqlite SqliteDB) checkUserExists(userid int64) (bool, error) {

	if err := sqlite.db.QueryRow(checkUsersQuery, userid).Scan(&userid); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}

//Wrapper function for insertToken
//Executes insertToken with a retry

func (sqlite SqliteDB) InsertToken(userid int64, token model.Token) error {
	var err error
	err = utils.Retry(func() error {
		err = sqlite.insertToken(userid, token)
		return err
	}, 2, time.Millisecond*20)
	return err
}

//Inserts the given token into the tokens table for the requested userid
//Returns error in case of failiure and nil in case of success

func (sqlite SqliteDB) insertToken(userid int64, token model.Token) error {
	//Insert token
	var err error
	if _, err = sqlite.db.Exec(insertNewTokenQuey, userid, token.Uuid, token.Expiration); err != nil {
		return err
	}

	return nil
}

//Wrapper function for getTokens
//Executes getTokens with a retry

func (sqlite SqliteDB) GetTokens(userid int64) ([]model.Token, error) {
	var getTokens []model.Token
	var err error
	err = utils.Retry(func() error {
		getTokens, err = sqlite.getTokens(userid)
		return err
	}, 2, time.Millisecond*20)
	return getTokens, err
}

//Recovers the access tokens for the requested userid (at most 2)
//Returns error in case of failiure and nil in case of success

func (sqlite SqliteDB) getTokens(userid int64) ([]model.Token, error) {

	var err error
	var rows *sql.Rows
	if rows, err = sqlite.db.Query(getFromTokensQuery, userid); err != nil {
		return nil, err
	}

	tokens := make([]model.Token, 0)
	defer rows.Close()
	for rows.Next() {
		token := model.Token{}
		err = rows.Scan(&token.Uuid, &token.Expiration)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	return tokens, nil
}

//Wrapper function for insertMessage
//Executes insertMessage with a retry

func (sqlite SqliteDB) InsertMessage(message model.MessageDTO) (model.MessageDTO, error) {
	var insertedMessage model.MessageDTO
	var err error
	err = utils.Retry(func() error {
		insertedMessage, err = sqlite.insertMessage(message)
		return err
	}, 2, time.Millisecond*20)
	return insertedMessage, err
}

//Inserts the given message into the messages table
//Returns error in case of failiure and nil in case of success

func (sqlite SqliteDB) insertMessage(message model.MessageDTO) (model.MessageDTO, error) {
	//Insert message
	var err error
	var insertResult sql.Result
	if insertResult, err = sqlite.db.Exec(insertNewMessageQuey, message.RecipientId, message.SenderId, message.Timestamp, message.Type, message.Text, message.Url, message.Height, message.Width, message.Source); err != nil {
		return message, err
	}
	id, err := insertResult.LastInsertId()
	message.MessageId = id
	return message, nil
}

//Wrapper function for getMessages
//Executes getMessages with a retry

func (sqlite SqliteDB) GetMessages(recipientId int64, messageId int64, limit int64) ([]model.MessageDTO, error) {
	var getMessages []model.MessageDTO
	var err error
	err = utils.Retry(func() error {
		getMessages, err = sqlite.getMessages(recipientId, messageId, limit)
		return err
	}, 2, time.Millisecond*20)
	return getMessages, err
}

//Recovers the access tokens for the requested userid (at most 2)
//Returns error in case of failiure and nil in case of success

func (sqlite SqliteDB) getMessages(recipientId int64, messageId int64, limit int64) ([]model.MessageDTO, error) {

	var err error
	var rows *sql.Rows
	if rows, err = sqlite.db.Query(getFromMessagesQuery, recipientId, messageId, limit); err != nil {
		return nil, err
	}

	messages := make([]model.MessageDTO, 0)
	defer rows.Close()
	for rows.Next() {
		message := model.MessageDTO{}
		err = rows.Scan(&message.MessageId, &message.RecipientId, &message.SenderId, &message.Timestamp, &message.Type, &message.Text, &message.Url, &message.Height, &message.Width, &message.Source)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}
