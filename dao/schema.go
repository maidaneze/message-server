package dao

const (
	//Gets the userid, password and salt from the users table.
	//The password is encrypted using RSASHA256 using the salt
	//The query is eficient because its performed on the INDEX "idx_username"

	getFromUsersQuery = "SELECT userid, password, password_salt FROM users WHERE username = ?"

	//Gets the userid from the users table for the given userid.

	checkUsersQuery = "SELECT userid FROM users WHERE userid = ?"

	//Inserts a user into the users table

	insertNewUserQuey = "INSERT INTO users (username, password, password_salt) VALUES (?, ?, ?)"

	//Inserts a new token into the tokens table.
	//Triggers the trigger "purgetokens" which ensures there are at most 2 tokens per user by deleting the oldest token.

	insertNewTokenQuey = "INSERT INTO tokens (userid, token, expiration) VALUES (?, ?, ?)"

	//Gets all the tokens corresponding to the user id (at most 2), which must then be verified

	getFromTokensQuery = "SELECT token, expiration FROM tokens WHERE userid = ?"

	//Inserts a new message into the messages table

	insertNewMessageQuey = "INSERT INTO messages (recipientid, senderid, timestamp, type, text, url, height, width, source) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"

	//Gets all the messages corresponding to the recipient id that have a messageid greater that the given messageid
	//Returns at most the given limit of rows

	getFromMessagesQuery = "SELECT messageid, recipientid, senderid, timestamp, type, text, url, height, width, source FROM messages WHERE recipientid = ? AND messageid >= ? ORDER BY messageid ASC LIMIT ?"

	//Users table schema
	//Allows efficent operations by using the index "idx_username" on the column username

	usersSchema = `CREATE TABLE users (userid INTEGER PRIMARY KEY, username TEXT, password TEXT, password_salt TEXT);
				CREATE UNIQUE INDEX idx_username ON users(username);`

	//Tokens table schema
	//Allows efficent operations by using the index "idx_tokens_userid" on the column userid
	//The "purgetoken" trigger ensures there are only a maximum of 2 active sessions at the time by deleting the oldest
	//token after each insert if there are more than 2 tokens in the database Sience its only triggered on updates,
	//it should be fairly efficent

	tokensSchema = `CREATE TABLE tokens (userid INTEGER,token TEXT,expiration INTEGER);
CREATE INDEX idx_tokens_userid ON tokens(userid);
CREATE TRIGGER purgetokens AFTER INSERT ON tokens
WHEN (SELECT count(*) FROM tokens WHERE userid = NEW.userid) > 2
BEGIN
DELETE FROM tokens WHERE userid = NEW.userid AND token in (SELECT token FROM tokens WHERE userid = NEW.userid ORDER BY expiration ASC LIMIT 1);
END;`

	//MessagesSchema
	//Allows efficent operations by using the index "idx_messages" first on the column reciever_userid and then on the column mesageid

	messagesSchema = `CREATE TABLE messages (
 messageid INTEGER PRIMARY KEY,
 recipientid INTEGER,
 senderid INTEGER,
 timestamp DATE,
 type TEXT CHECK ( type IN ('text','image','video')),
 text TEXT,
 url TEXT,
 height INTEGER,
 width INTEGER,
 source TEXT CHECK ( source IN ('youtube','vimeo',''))
);
CREATE INDEX idx_messages ON messages(recipientid, messageid);`
)
