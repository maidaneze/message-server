package dao

import (
	"github.com/maidaneze/message-server/utils"
	"os"
	"testing"
	"time"
)

//Recreates the database
func RefreshSchema(sqlite SqliteDB) {
	refreshUsersSchema(sqlite)
	refreshTokensSchema(sqlite)
	refreshMessagesSchema(sqlite)
}

//Creates the users table
//Returns error if it fails and nil if it succeeds

func createUsersSchema(sqlite SqliteDB) error {
	sqlStmt := usersSchema

	var err error
	utils.Retry(func() error {
		_, err = sqlite.db.Exec(sqlStmt)
		return err
	}, 2, time.Millisecond*20)
	return err
}

//Recreates the users table
//Returns error if it fails and nil if it succeeds

func refreshUsersSchema(sqlite SqliteDB) error {
	sqlStmt := `drop table users;`

	var err error
	utils.Retry(func() error {
		_, err = sqlite.db.Exec(sqlStmt)
		return err
	}, 2, time.Millisecond*20)

	return createUsersSchema(sqlite)
}

//Recreates the tokens table
//Returns error if it fails and nil if it succeeds

func createTokensSchema(sqlite SqliteDB) error {
	sqlStmt := tokensSchema

	var err error
	utils.Retry(func() error {
		_, err = sqlite.db.Exec(sqlStmt)
		return err
	}, 2, time.Millisecond*20)
	return err
}

//Recreates the tokens table
//Returns error if it fails and nil if it succeeds

func refreshTokensSchema(sqlite SqliteDB) error {
	sqlStmt := `drop table tokens;`

	var err error
	utils.Retry(func() error {
		_, err = sqlite.db.Exec(sqlStmt)
		return err
	}, 2, time.Millisecond*20)

	return createTokensSchema(sqlite)
}

//Recreates the messages table
//Returns error if it fails and nil if it succeeds

func createMessagesSchema(sqlite SqliteDB) error {
	sqlStmt := messagesSchema

	var err error
	utils.Retry(func() error {
		_, err = sqlite.db.Exec(sqlStmt)
		return err
	}, 2, time.Millisecond*20)
	return err
}

//Recreates the messages table
//Returns error if it fails and nil if it succeeds

func refreshMessagesSchema(sqlite SqliteDB) error {
	sqlStmt := `drop table messages;`

	var err error
	utils.Retry(func() error {
		_, err = sqlite.db.Exec(sqlStmt)
		return err
	}, 2, time.Millisecond*20)

	return createMessagesSchema(sqlite)
}

func SetupSqliteDatabaseTest(t *testing.T, filename string) SqliteDB {
	os.Remove(filename)

	db, err := OpenSqlite3Database(filename)

	if err != nil {
		t.FailNow()
	}

	RefreshSchema(db)
	return db
}

func TeardownSqliteDatabaseTest(sqlite SqliteDB, filename string) {
	sqlite.db.Close()
	os.Remove(filename)
}
