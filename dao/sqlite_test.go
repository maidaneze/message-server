package dao

import (
	"math"
	"os"
	"testing"

	"strings"

	"github.com/maidaneze/message-server/services/users"

	"github.com/maidaneze/message-server/model"
	"github.com/maidaneze/message-server/services/auth"

	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testDatabase          SqliteDB
	testDatabaseFiletName string = "./foo.db"
)

//Test suites

func TestSqlite3DatabaseSuiteWithOpenConnection(t *testing.T) {
	//Setup

	testDatabase = SetupSqliteDatabaseTest(t, testDatabaseFiletName)

	//TestSuite

	t.Run("testCheckConnectionShouldSucceedIfConnectionIsOpen", testCheckConnectionShouldSucceedIfConnectionIsOpen)
	t.Run("testInsertUserShouldSaveUsersProperly", testInsertUserShouldSaveUsersProperly)
	t.Run("testGetUserNotStoredShouldReturnFalseAndNoError", testGetUserNotStoredShouldReturnFalseAndNoError)
	t.Run("testCheckUserExistsNotStoredShouldReturnFalseAndNoError", testCheckUserExistsNotStoredShouldReturnFalseAndNoError)
	t.Run("testInsertRepeatedUserShouldFail", testInsertRepeatedUserShouldFail)
	t.Run("testInsertRepeatedUserShouldReturnIncreasingIds", testInsertRepeatedUserShouldReturnIncreasingIds)
	t.Run("testInsertTokenShouldSaveTokensProperly", testInsertTokenShouldSaveTokensProperly)
	t.Run("testInsertdTokensShouldBePurged", testInsertdTokensShouldBePurged)
	t.Run("testInsertMessageShouldSaveMessageProperly", testInsertMessageShouldSaveMessageProperly)
	t.Run("testInsertMessageShouldFailOnInvalidTypeOrSource", testInsertMessageShouldFailOnInvalidTypeOrSource)
	t.Run("testGetMessagesShouldGetAsMuchAsLimitNumberOfMessages", testGetMessagesShouldGetAsMuchAsLimitNumberOfMessages)
	t.Run("testGetMessagesShouldGetMessagesStartingFromMessageId", testGetMessagesShouldGetMessagesStartingFromMessageId)

	//Teardown

	TeardownSqliteDatabaseTest(testDatabase, testDatabaseFiletName)
}

func TestSqlite3DatabaseSuiteWithClosedConnection(t *testing.T) {
	//Setup
	testDatabase = SetupSqliteDatabaseTest(t, testDatabaseFiletName)
	testDatabase.db.Close()

	//TestSuite

	t.Run("testCheckConnectionShouldFailIfConectionIsNotOpen", testCheckConnectionShouldFailIfConectionIsNotOpen)
	t.Run("testGetUserShouldFailWithClosedConnection", testGetUserShouldFailWithClosedConnection)
	t.Run("testCheckUserExistsShouldFailWithClosedConnection", testCheckUserExistsShouldFailWithClosedConnection)
	t.Run("testInsertUserShouldFailWithClosedConnection", testInsertUserShouldFailWithClosedConnection)
	t.Run("testInsertTokenShouldFailWithClosedConnection", testInsertTokenShouldFailWithClosedConnection)
	t.Run("testGetTokensShouldFailWithClosedConnection", testGetTokensShouldFailWithClosedConnection)
	t.Run("testInsertMessageShouldFailWithClosedConnection", testInsertMessageShouldFailWithClosedConnection)
	t.Run("testGetMessagesShouldFailWithClosedConnection", testGetMessagesShouldFailWithClosedConnection)

	//TearDown
	os.Remove(testDatabaseFiletName)
}

//Actual tests

func testCheckConnectionShouldFailIfConectionIsNotOpen(t *testing.T) {
	err := testDatabase.CheckConnection()
	assert.NotNil(t, err)
}

func testCheckConnectionShouldSucceedIfConnectionIsOpen(t *testing.T) {
	RefreshSchema(testDatabase)
	err := testDatabase.CheckConnection()
	assert.Nil(t, err)
}

func testGetUserShouldFailWithClosedConnection(t *testing.T) {
	cases := []struct {
		name     string
		username string
	}{
		{"getUserShouldFailTestWithEmptyFieltds", ""},
		{"getUserShouldFailTestSampleUser", "user1"},
		{"getUserShouldFailRepeatTestCase", "user1"},
		{"getUserShouldFailTestSampleUserWithNumbers", "123"},
		{"getUserShouldFailTestSampleUserWeirdCharactersInUser", "$%·WeirdCharacters=)("},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			_, _, err := testDatabase.GetUser(c.username)
			assert.NotNil(tt, err)
		})
	}
}

func testCheckUserExistsShouldFailWithClosedConnection(t *testing.T) {
	cases := []struct {
		name   string
		userid int64
	}{
		{"checkUserExistsShouldFailTestWithEmptyFieltds0", 0},
		{"checkUserExistsShouldFailTestWithEmptyFieltds1", 1},
		{"checkUserExistsShouldFailTestWithEmptyFieltds2", 2},
		{"checkUserExistsShouldFailTestWithEmptyFieltds3", 3},
		{"checkUserExistsShouldFailTestWithEmptyFieltds4", 4},
		{"checkUserExistsShouldFailTestWithEmptyFieltdsMaxint", math.MaxInt64},
		{"checkUserExistsShouldFailTestWithEmptyFieltds1", -1},
		{"checkUserExistsShouldFailTestWithEmptyFieltds2", -2},
		{"checkUserExistsShouldFailTestWithEmptyFieltds3", -3},
		{"checkUserExistsShouldFailTestWithEmptyFieltds4", -4},
		{"checkUserExistsShouldFailTestWithEmptyFieltdsMinint", math.MinInt64},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			_, err := testDatabase.CheckUserExists(c.userid)
			assert.NotNil(t, err)
		})
	}
}

func testInsertUserShouldFailWithClosedConnection(t *testing.T) {
	cases := []struct {
		name     string
		username string
		password string
	}{
		{"testInsertShouldFailTestWithEmptyFieltds", "", ""},
		{"testInsertShouldFailTestSampleUser", "user1", "123"},
		{"testInsertShouldFailTepeatTestCase", "user1", "123"},
		{"testInsertShouldFailTestSampleUserStringPass", "user2", "notNumbers"},
		{"testInsertShouldFailTestSampleUserWeirdCharactersInPass", "user3", "$%·WeirdCharacters=)("},
		{"testInsertShouldFailTestSampleUserMadeOfNumebrs", "123", "pass1"},
		{"testInsertShouldFailTestSampleUserWeirdCharactersInUser", "$%·WeirdCharacters=)(", "pass1"},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			user, err := users.CreateUser(c.username, c.password)
			_, err = testDatabase.InsertUser(user)
			assert.NotNil(tt, err)
		})
	}
}

func testInsertUserShouldSaveUsersProperly(t *testing.T) {
	RefreshSchema(testDatabase)
	cases := []struct {
		name     string
		username string
		password string
		setup    func(db SqliteDB)
	}{
		{"createUserShouldSucceedWithEmptyFieltds", "", "", RefreshSchema},
		{"createUserShouldSucceedTestSampleUser", "user1", "123", RefreshSchema},
		{"createUserShouldSucceedTestSampleUserStringPass", "user2", "notNumbers", RefreshSchema},
		{"createUserShouldSucceedTestSampleUserWeirdCharactersInPass", "user3", "$%·WeirdCharacters=)(", RefreshSchema},
		{"createUserShouldSucceedTestSampleUserMadeOfNumebrs", "123", "pass1", RefreshSchema},
		{"createUserShouldSucceedTestSampleUserWeirdCharactersInUser", "$%·WeirdCharacters=)(", "pass1", RefreshSchema},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {

			user, err := users.CreateUser(c.username, c.password)
			createdUser, err := testDatabase.InsertUser(user)

			assert.Nil(tt, err)

			var found bool
			found, err = testDatabase.CheckUserExists(createdUser.Userid)

			assert.True(tt, found)
			assert.Nil(tt, err)

			user, found, err = testDatabase.GetUser(c.username)

			assert.True(tt, found)
			assert.Nil(tt, err)
			assert.Equal(tt, createdUser.Userid, user.Userid)
			assert.Equal(tt, createdUser.Password, user.Password)
		})
	}
}

func testGetUserNotStoredShouldReturnFalseAndNoError(t *testing.T) {
	RefreshSchema(testDatabase)
	username := "username"

	_, found, err := testDatabase.GetUser(username)

	assert.False(t, found)
	assert.Nil(t, err)
}

func testCheckUserExistsNotStoredShouldReturnFalseAndNoError(t *testing.T) {
	RefreshSchema(testDatabase)
	userid := int64(0)

	found, err := testDatabase.checkUserExists(userid)

	assert.False(t, found)
	assert.Nil(t, err)
}

func testInsertRepeatedUserShouldFail(t *testing.T) {
	RefreshSchema(testDatabase)

	username := "user1"
	pass1 := "123"
	pass2 := "456"
	user, err := users.CreateUser(username, pass1)
	_, err = testDatabase.InsertUser(user)
	assert.Nil(t, err)

	cases := []struct {
		name     string
		password string
	}{
		{"createUserShouldFailTestWhenUsingRepeatedUsernameAndPassword", pass1},
		{"createUserShouldFailTestWhenUsingRepeatedUsernameAndDifferentPassword", pass2},
	}
	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			user, err := users.CreateUser(username, c.password)
			_, err = testDatabase.InsertUser(user)
			assert.NotNil(tt, err)
		})
	}
}

func testInsertRepeatedUserShouldReturnIncreasingIds(t *testing.T) {
	RefreshSchema(testDatabase)
	username := "user"
	pass1 := "123"

	user, err := users.CreateUser(username+"1", pass1)
	user1, err := testDatabase.InsertUser(user)
	assert.Nil(t, err)
	user, err = users.CreateUser(username+"2", pass1)
	assert.Nil(t, err)
	user2, err := testDatabase.InsertUser(user)
	assert.Nil(t, err)
	user, err = users.CreateUser(username+"3", pass1)
	assert.Nil(t, err)
	user3, err := testDatabase.InsertUser(user)
	assert.Nil(t, err)

	assert.Equal(t, int64(1), user1.Userid)
	assert.Equal(t, int64(2), user2.Userid)
	assert.Equal(t, int64(3), user3.Userid)
}

//Helper functions

func testInsertTokenShouldFailWithClosedConnection(t *testing.T) {
	cases := []struct {
		name   string
		userid int64
	}{
		{"insertTokenShouldFailTestWithZeroValue", 0},
		{"insertTokenShouldFailTestSampleUser", 123},
		{"insertTokenShouldFailTestSampleUser2", 456},
		{"insertTokenShouldFailTestSampleUser3", 789},
		{"insertTokenShouldFailTepeatTestCase", 123},
		{"insertTokenShouldFailWithBigNumber", math.MaxInt64},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			err := testDatabase.InsertToken(c.userid, model.Token{})
			assert.NotNil(tt, err)
		})
	}
}
func testGetTokensShouldFailWithClosedConnection(t *testing.T) {
	cases := []struct {
		name   string
		userid int64
	}{
		{"getTokenShouldFailTestWithZeroValue", 0},
		{"getTokenShouldFailTestSampleUser", 123},
		{"getTokenShouldFailTestSampleUser2", 456},
		{"getTokenShouldFailTestSampleUser3", 789},
		{"getTokenShouldFailTepeatTestCase", 123},
		{"getTokenShouldFailWithBigNumber", math.MaxInt64},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			_, err := testDatabase.GetTokens(c.userid)
			assert.NotNil(tt, err)
		})
	}
}

func testInsertMessageShouldFailWithClosedConnection(t *testing.T) {
	cases := []struct {
		name        string
		recipientid int64
	}{
		{"insertMessageShouldFailTestWithZeroValue", 0},
		{"insertMessageShouldFailTestSampleUser", 123},
		{"insertMessageShouldFailTestSampleUser2", 456},
		{"insertMessageShouldFailTestSampleUser3", 789},
		{"insertMessageShouldFailTepeatTestCase", 123},
		{"insertMessageShouldFailWithBigNumber", math.MaxInt64},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			message := model.MessageDTO{}
			message.RecipientId = c.recipientid
			_, err := testDatabase.InsertMessage(message)
			assert.NotNil(tt, err)
		})
	}
}
func testGetMessagesShouldFailWithClosedConnection(t *testing.T) {
	cases := []struct {
		name        string
		recipientid int64
	}{
		{"getMessagesShouldFailTestWithZeroValue", 0},
		{"getMessagesShouldFailTestSampleUser", 123},
		{"getMessagesShouldFailTestSampleUser2", 456},
		{"getMessagesShouldFailTestSampleUser3", 789},
		{"getMessagesShouldFailTepeatTestCase", 123},
		{"getMessagesShouldFailWithBigNumber", math.MaxInt64},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			_, err := testDatabase.GetMessages(c.recipientid, 0, 100)
			assert.NotNil(tt, err)
		})
	}
}

func testInsertTokenShouldSaveTokensProperly(t *testing.T) {
	RefreshSchema(testDatabase)
	cases := []struct {
		name   string
		userid int64
	}{
		{"insertTokenShouldSucceedTestWithZeroValue", 0},
		{"insertTokenShouldSucceedTestSampleUser", 123},
		{"insertTokenShouldSucceedTestSampleUser2", 456},
		{"insertTokenShouldSucceedTestSampleUser3", 789},
		{"insertTokenShouldSucceedTepeatTestCase", 123},
		{"insertTokenShouldSucceedWithBigNumber", math.MaxInt64},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			token, _ := auth.GenerateToken()
			err := testDatabase.InsertToken(c.userid, token)
			assert.Nil(tt, err)
			tokens, err := testDatabase.GetTokens(c.userid)
			assert.Nil(tt, err)
			assert.True(tt, auth.FindToken(token, tokens))
		})
	}
}

func testInsertMessageShouldSaveMessageProperly(t *testing.T) {
	RefreshSchema(testDatabase)
	cases := []struct {
		name              string
		recipientid       int64
		senderid          int64
		messageType       string
		text              string
		url               string
		height            int64
		width             int64
		source            string
		expectedMessageId int64
	}{
		{"insertMessageShouldSucceedTestWithZeroValue1", 0, 0, "text", "", "", 0, 0, "youtube", 1},
		{"insertMessageShouldSucceedTestWithZeroValue2", 0, 0, "image", "", "", 0, 0, "youtube", 2},
		{"insertMessageShouldSucceedTestWithZeroValue3", 0, 0, "video", "", "", 0, 0, "youtube", 3},
		{"insertMessageShouldSucceedTestWithZeroValue4", 0, 0, "text", "", "", 0, 0, "vimeo", 4},
		{"insertMessageShouldSucceedTestWithZeroValue5", 0, 0, "image", "", "", 0, 0, "vimeo", 5},
		{"insertMessageShouldSucceedTestWithZeroValue6", 0, 0, "video", "", "", 0, 0, "vimeo", 6},
		{"insertMessageShouldSucceedTestWithRecipientId", 1, 0, "text", "", "", 0, 0, "vimeo", 7},
		{"insertMessageShouldSucceedTestWithSenderId", 0, 1, "text", "", "", 0, 0, "vimeo", 8},
		{"insertMessageShouldSucceedTestWithText", 0, 0, "text", "text", "", 0, 0, "vimeo", 9},
		{"insertMessageShouldSucceedTestWithUrl", 0, 0, "text", "", "url", 0, 0, "youtube", 10},
		{"insertMessageShouldSucceedTestWithHeight", 0, 0, "text", "", "", 150, 0, "youtube", 11},
		{"insertMessageShouldSucceedTestWithWidth", 0, 0, "text", "", "", 0, 150, "youtube", 12},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			insertMessage := model.MessageDTO{0, c.recipientid, c.senderid, time.Now(), c.messageType, c.text, c.url, c.height, c.width, c.source}
			message, err := testDatabase.InsertMessage(insertMessage)
			assert.Nil(tt, err)

			assert.Equal(tt, c.expectedMessageId, message.MessageId)
			assert.Equal(tt, c.recipientid, message.RecipientId)
			assert.Equal(tt, c.senderid, message.SenderId)
			assert.Equal(tt, c.messageType, message.Type)
			assert.Equal(tt, c.text, message.Text)
			assert.Equal(tt, c.url, message.Url)
			assert.Equal(tt, c.height, message.Height)
			assert.Equal(tt, c.width, message.Width)
			assert.Equal(tt, c.source, message.Source)

			messages, err := testDatabase.GetMessages(c.recipientid, message.MessageId, 1)
			require.True(tt, len(messages) == 1)
			getMessage := messages[0]
			assert.Nil(tt, err)

			assert.Equal(tt, c.expectedMessageId, getMessage.MessageId)
			assert.Equal(tt, c.recipientid, getMessage.RecipientId)
			assert.Equal(tt, c.senderid, getMessage.SenderId)
			assert.Equal(tt, c.messageType, getMessage.Type)
			assert.Equal(tt, c.text, getMessage.Text)
			assert.Equal(tt, c.url, getMessage.Url)
			assert.Equal(tt, c.height, getMessage.Height)
			assert.Equal(tt, c.width, getMessage.Width)
			assert.Equal(tt, c.source, getMessage.Source)
		})
	}
}

func testInsertMessageShouldFailOnInvalidTypeOrSource(t *testing.T) {
	RefreshSchema(testDatabase)
	cases := []struct {
		name        string
		messageType string
		source      string
		success     bool
	}{
		{"insertMessageShouldSucceedWithValidSourceAndType1", "text", "youtube", true},
		{"insertMessageShouldSucceedWithValidSourceAndType2", "image", "youtube", true},
		{"insertMessageShouldSucceedWithValidSourceAndType3", "video", "youtube", true},
		{"insertMessageShouldSucceedWithValidSourceAndType4", "text", "vimeo", true},
		{"insertMessageShouldSucceedWithValidSourceAndType5", "image", "vimeo", true},
		{"insertMessageShouldSucceedWithValidSourceAndType6", "video", "vimeo", true},
		{"insertMessageShouldSucceedWithInvalidSourceAndValidType1", "text", "invalid", false},
		{"insertMessageShouldSucceedWithInvalidSourceAndValidType2", "image", "invalid", false},
		{"insertMessageShouldSucceedWithInvalidSourceAndValidType3", "video", "invalid", false},
		{"insertMessageShouldSucceedWithValidSourceAndInvalidType1", "invalid", "youtube", false},
		{"insertMessageShouldSucceedWithValidSourceAndInvalidType2", "invalid", "vimeo", false},
		{"insertMessageShouldSucceedWithInvalidSourceAndInvalidType", "invalid", "invalid", false},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			insertMessage := model.MessageDTO{0, 0, 0, time.Now(), c.messageType, "", "", 0, 0, c.source}
			_, err := testDatabase.InsertMessage(insertMessage)
			assert.True(tt, c.success == (err == nil))
		})
	}
}

func testGetMessagesShouldGetAsMuchAsLimitNumberOfMessages(t *testing.T) {
	RefreshSchema(testDatabase)
	cases := []struct {
		name             string
		saveMessages     int
		limit            int64
		expectedMessages int
	}{
		{"getMessageShouldReturnLessThanLimit1", 0, 5, 0},
		{"getMessageShouldReturnLessThanLimit1", 1, 5, 1},
		{"getMessageShouldReturnLessThanLimit1", 2, 5, 2},
		{"getMessageShouldReturnLessThanLimit1", 3, 5, 3},
		{"getMessageShouldReturnLessThanLimit1", 4, 5, 4},
		{"getMessageShouldReturnLimit1", 5, 5, 5},
		{"getMessageShouldReturnLimit2", 6, 5, 5},
		{"getMessageShouldReturnLimit3", 7, 5, 5},
		{"getMessageShouldReturnLimit4", 8, 5, 5},
		{"getMessageShouldReturnLimit5", 9, 5, 5},
		{"getMessageShouldReturnLimit6", 10, 5, 5},
		{"getMessageShouldReturnLimit7", 10, 6, 6},
		{"getMessageShouldReturnLimit8", 10, 7, 7},
		{"getMessageShouldReturnLimit9", 10, 8, 8},
		{"getMessageShouldReturnLimit10", 10, 9, 9},
		{"getMessageShouldReturnLimit11", 10, 10, 10},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			var id int64 = 0
			for i := 0; i < c.saveMessages; i++ {
				insertMessage := model.MessageDTO{0, 0, 0, time.Now(), "text", "", "", 0, 0, "youtube"}
				message, err := testDatabase.InsertMessage(insertMessage)
				assert.Nil(tt, err)
				if i == 0 {
					id = message.MessageId
				}

			}

			messages, err := testDatabase.GetMessages(0, id, c.limit)
			assert.Nil(tt, err)
			assert.Equal(tt, c.expectedMessages, len(messages))
		})
	}
}

func testGetMessagesShouldGetMessagesStartingFromMessageId(t *testing.T) {
	RefreshSchema(testDatabase)
	var id int64 = 0
	for i := 0; i < 5; i++ {
		insertMessage := model.MessageDTO{0, 0, 0, time.Now(), "text", "", "", 0, 0, "youtube"}
		message, err := testDatabase.InsertMessage(insertMessage)
		assert.Nil(t, err)
		if i == 0 {
			id = message.MessageId
		}

	}

	for i := 0; i < 5; i++ {
		messages, err := testDatabase.GetMessages(0, id, 5)
		assert.Nil(t, err)
		for j := 0; j < len(messages); j++ {
			assert.Equal(t, id+int64(j), messages[j].MessageId)
		}
	}
}

func testInsertdTokensShouldBePurged(t *testing.T) {
	RefreshSchema(testDatabase)

	var userid1 int64 = 123
	var userid2 int64 = 456

	token1, _ := auth.GenerateToken()
	token2, _ := auth.GenerateToken()
	token3, _ := auth.GenerateToken()
	token4, _ := auth.GenerateToken()
	token5, _ := auth.GenerateToken()
	err := testDatabase.InsertToken(userid1, token1)
	assert.Nil(t, err)

	err = testDatabase.InsertToken(userid1, token2)
	assert.Nil(t, err)

	err = testDatabase.InsertToken(userid2, token3)
	assert.Nil(t, err)

	err = testDatabase.InsertToken(userid2, token4)
	assert.Nil(t, err)

	err = testDatabase.InsertToken(userid2, token5)
	assert.Nil(t, err)

	tokens1, err := testDatabase.GetTokens(userid1)
	assert.Nil(t, err)

	tokens2, err := testDatabase.GetTokens(userid2)
	assert.Nil(t, err)

	//Ensures the query brings all elements that it finds
	assert.False(t, strings.Contains(getFromTokensQuery, "LIMIT"))

	assert.Equal(t, 2, len(tokens1))
	assert.Equal(t, 2, len(tokens2))

	assert.True(t, auth.FindToken(token1, tokens1))
	assert.False(t, auth.FindToken(token1, tokens2))

	assert.True(t, auth.FindToken(token2, tokens1))
	assert.False(t, auth.FindToken(token2, tokens2))

	assert.False(t, auth.FindToken(token3, tokens1))
	assert.False(t, auth.FindToken(token3, tokens2))

	assert.False(t, auth.FindToken(token4, tokens1))
	assert.True(t, auth.FindToken(token4, tokens2))

	assert.False(t, auth.FindToken(token5, tokens1))
	assert.True(t, auth.FindToken(token5, tokens2))

}
