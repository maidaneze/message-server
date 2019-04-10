package controllers

import (
	"github.com/maidaneze/message-server/dao"
	"net/http"
	"testing"

	"encoding/json"
	"io/ioutil"

	"bytes"

	"github.com/maidaneze/message-server/model"

	"github.com/maidaneze/message-server/services/auth"

	"math"

	"fmt"

	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testHandler Handler
	mockServer  *http.Server
	testDB      dao.SqliteDB
)

func TestSqlite3DatabaseSuiteWithOpenConnection(t *testing.T) {
	//Setup

	testDB = dao.SetupSqliteDatabaseTest(t, "foo.db")
	testHandler = Handler{testDB}

	mockServer = testHandler.Setup()
	go mockServer.ListenAndServe()

	//TestSuiteWithOpenDB

	t.Run("testHandler_CheckReturnsOk", testHandler_CheckReturnsOk)
	t.Run("testRegisterUser", testRegisterUser)
	t.Run("testFailToRegisterUserTwice", testFailToRegisterUserTwice)
	t.Run("testFailToRegisterUserInvalidFields", testFailToRegisterUserInvalidFields)
	t.Run("testFailToCreateUserInvalidBody", testFailToCreateUserInvalidBody)
	t.Run("testLoginUser", testLoginUser)
	t.Run("testFailToLoginUserInvalidBody", testFailToLoginUserInvalidBody)
	t.Run("testFailToLoginUnregisteredUser", testFailToLoginUnregisteredUser)
	t.Run("testFailToLoginWrongPassword", testFailToLoginWrongPassword)
	t.Run("testFailToLoginUserInvalidFields", testFailToLoginUserInvalidFields)
	t.Run("testFailToPostTextMessageInvalidFields", testFailToPostTextMessageInvalidFields)
	t.Run("testFailToPostImageMessageInvalidFields", testFailToPostImageMessageInvalidFields)
	t.Run("testFailToPostVideoMessageInvalidFields", testFailToPostVideoMessageInvalidFields)
	t.Run("testFailToPostMessageInvalidBody", testFailToPostMessageInvalidBody)
	t.Run("testFailToPostMessageUnauthorized", testFailToPostMessageUnauthorized)
	t.Run("testPostTextMessage", testPostTextMessage)
	t.Run("testPostImageMessage", testPostImageMessage)
	t.Run("testPostVideoMessage", testPostVideoMessage)
	t.Run("testFailToGetMessageInvalidFields", testFailToGetMessageInvalidFields)
	t.Run("testFailToGetMessageInvalidQueryParams", testFailToGetMessageInvalidQueryParams)
	t.Run("testFailToGetMessageUnauthorized", testFailToGetMessageUnauthorized)
	t.Run("testGetMessages", testGetMessages)
	//Test de exito

	//CloseDB

	dao.TeardownSqliteDatabaseTest(testDB, "foo.db")

	//TestSuiteWithClosedDB

	t.Run("testHandler_CheckReturnsInternalServerError", testHandler_CheckReturnsInternalServerError)
	t.Run("testFailToRegisterUserWithClosedConnection", testFailToRegisterUserWithClosedConnection)
	t.Run("testFailToLoginUserWithClosedConnection", testFailToLoginUserWithClosedConnection)
	t.Run("testFailToPostMessageWithClosedConnection", testFailToPostMessageWithClosedConnection)
	t.Run("testFailToGetMessageWithClosedConnection", testFailToGetMessageWithClosedConnection)

	//Teardown

	mockServer.Shutdown(nil)
}

func testHandler_CheckReturnsOk(t *testing.T) {
	resp, err := http.Post("http://localhost:8080/check", "application/json", nil)
	require.Nil(t, err)
	defer resp.Body.Close()

	assert.Equal(t, resp.StatusCode, http.StatusOK)

	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)

	response := struct {
		Health string `json:"health"`
	}{}
	json.Unmarshal(body, &response)

	assert.Equal(t, "ok", response.Health)
}

func testHandler_CheckReturnsInternalServerError(t *testing.T) {
	resp, err := http.Post("http://localhost:8080/check", "application/json", nil)
	require.Nil(t, err)
	defer resp.Body.Close()

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
}

func testRegisterUser(t *testing.T) {
	dao.RefreshSchema(testDB)
	nextId := int64(1)
	cases := []struct {
		name     string
		username string
		password string
	}{
		{"testRegisterUser1", "user1", "pass1"},
		{"testRegisterUser2", "user2", "pass2"},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			id := createUserSuccessfullyForTest(tt, c.username, c.password)
			assert.Equal(tt, nextId, id)
			nextId++
		})
	}
}

func testFailToRegisterUserInvalidFields(t *testing.T) {
	dao.RefreshSchema(testDB)
	bigUser := ""
	for i := 0; i < model.MAX_USERNAME_FIELD_SIZE; i += 10 {
		bigUser = bigUser + "hugestring"
	}
	bigPassword := ""
	for i := 0; i < model.MAX_PASSWORD_FIELD_SIZE; i += 10 {
		bigPassword = bigPassword + "hugestring"
	}
	cases := []struct {
		name     string
		username string
		password string
	}{
		{"testFailToRegisterUserInvalidBodyNoUser", "", "pass1"},
		{"testFailToRegisterUserInvalidBodyPassword", "user1", ""},
		{"testFailToRegisterUserInvalidBodyBigUser", bigUser, "pass1"},
		{"testFailToRegisterUserInvalidBodyBigPassword", "user1", bigPassword},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			resp, err := requestCreateUser(c.username, c.password)
			if err == nil {
				defer resp.Body.Close()
			}
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}
}

func testFailToRegisterUserTwice(t *testing.T) {
	dao.RefreshSchema(testDB)

	username := "user1"
	password := "pass1"
	_ = createUserSuccessfullyForTest(t, username, password)
	resp, err := requestCreateUser(username, password)
	if err == nil {
		defer resp.Body.Close()
	}
	assert.Equal(t, http.StatusConflict, resp.StatusCode)
}

func testFailToCreateUserInvalidBody(t *testing.T) {
	dao.RefreshSchema(testDB)

	requestReader := bytes.NewReader([]byte("invalid"))

	resp, _ := http.Post("http://localhost:8080/users", "application/json", requestReader)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func testFailToRegisterUserWithClosedConnection(t *testing.T) {
	dao.RefreshSchema(testDB)

	username := "user1"
	password := "pass1"
	resp, err := requestCreateUser(username, password)
	if err == nil {
		defer resp.Body.Close()
	}
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func testLoginUser(t *testing.T) {
	dao.RefreshSchema(testDB)
	cases := []struct {
		name     string
		username string
		password string
	}{
		{"testLogin1", "user1", "pass1"},
		{"testLogin2", "user2", "pass2"},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			_ = createUserSuccessfullyForTest(tt, c.username, c.password)
			_, _ = loginSuccessfullyForTest(t, c.username, c.password)
		})
	}
}

func testFailToLoginUserInvalidFields(t *testing.T) {
	dao.RefreshSchema(testDB)
	bigUser := ""
	for i := 0; i < model.MAX_USERNAME_FIELD_SIZE; i += 10 {
		bigUser = bigUser + "hugestring"
	}
	bigPassword := ""
	for i := 0; i < model.MAX_PASSWORD_FIELD_SIZE; i += 10 {
		bigPassword = bigPassword + "hugestring"
	}
	cases := []struct {
		name     string
		username string
		password string
	}{
		{"testFailToLoginUserInvalidBodyNoUser", "", "pass1"},
		{"testFailToLoginUserInvalidBodyPassword", "user1", ""},
		{"testFailToLoginUserInvalidBodyBigUser", bigUser, "pass1"},
		{"testFailToLoginUserInvalidBodyBigPassword", "user1", bigPassword},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			resp, err := requestCreateUser(c.username, c.password)
			if err == nil {
				defer resp.Body.Close()
			}
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}
}

func testFailToLoginUnregisteredUser(t *testing.T) {
	dao.RefreshSchema(testDB)

	username := "user1"
	password := "pass1"
	resp, err := requestLogin(username, password)
	if err == nil {
		defer resp.Body.Close()
	}
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func testFailToLoginWrongPassword(t *testing.T) {
	dao.RefreshSchema(testDB)

	username := "user1"
	password := "pass1"
	_ = createUserSuccessfullyForTest(t, username, password)
	resp, err := requestLogin(username, "invalid")
	if err == nil {
		defer resp.Body.Close()
	}
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func testFailToLoginUserInvalidBody(t *testing.T) {
	dao.RefreshSchema(testDB)

	requestReader := bytes.NewReader([]byte("invalid"))

	resp, _ := http.Post("http://localhost:8080/login", "application/json", requestReader)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func testFailToLoginUserWithClosedConnection(t *testing.T) {
	dao.RefreshSchema(testDB)

	username := "user1"
	password := "pass1"
	resp, err := requestLogin(username, password)
	if err == nil {
		defer resp.Body.Close()
	}
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func testPostTextMessage(t *testing.T) {
	dao.RefreshSchema(testDB)
	username1 := "user1"
	password1 := "pass1"

	username2 := "user2"
	password2 := "pass2"
	id1 := createUserSuccessfullyForTest(t, username1, password1)
	_ = createUserSuccessfullyForTest(t, username2, password2)
	id2, token2 := loginSuccessfullyForTest(t, username2, password2)
	content := model.TextContent{"text", "text"}
	id, _ := postMessageSuccessfullyForTest(t, id2, id1, content, token2)
	assert.Equal(t, int64(1), id)

}

func testPostImageMessage(t *testing.T) {
	dao.RefreshSchema(testDB)
	username1 := "user1"
	password1 := "pass1"

	username2 := "user2"
	password2 := "pass2"
	id1 := createUserSuccessfullyForTest(t, username1, password1)
	_ = createUserSuccessfullyForTest(t, username2, password2)
	id2, token2 := loginSuccessfullyForTest(t, username2, password2)
	content := model.ImageContent{"image", "www.image.com", int64(150), int64(120)}
	id, _ := postMessageSuccessfullyForTest(t, id2, id1, content, token2)
	assert.Equal(t, int64(1), id)
}

func testPostVideoMessage(t *testing.T) {
	dao.RefreshSchema(testDB)
	username1 := "user1"
	password1 := "pass1"

	username2 := "user2"
	password2 := "pass2"
	id1 := createUserSuccessfullyForTest(t, username1, password1)
	_ = createUserSuccessfullyForTest(t, username2, password2)
	id2, token2 := loginSuccessfullyForTest(t, username2, password2)
	content := model.VideoContent{"video", "www.youtube.com", "youtube"}
	id, _ := postMessageSuccessfullyForTest(t, id2, id1, content, token2)
	assert.Equal(t, int64(1), id)
}

func testFailToPostTextMessageInvalidFields(t *testing.T) {
	dao.RefreshSchema(testDB)

	bigTextField := ""
	for i := 0; i < model.MAX_TEXT_FIELD_SIZE; i += 10 {
		bigTextField = bigTextField + "hugestring"
	}

	username1 := "user1"
	password1 := "pass1"

	username2 := "user2"
	password2 := "pass2"
	id1 := createUserSuccessfullyForTest(t, username1, password1)
	_ = createUserSuccessfullyForTest(t, username2, password2)
	id2, token2 := loginSuccessfullyForTest(t, username2, password2)
	content := model.TextContent{"text", "text"}
	bigContent := model.TextContent{"text", bigTextField}

	cases := []struct {
		name        string
		senderid    int64
		recipientid int64
		content     model.MessageContent
	}{
		{"testFailToPostTextMessageSenderZero", 0, id1, content},
		{"testFailToPostTextMessageRecipientZero", id2, 0, content},
		{"testFailToPostTextMessageUnregisteredSender", math.MaxInt64, id1, content},
		{"testFailToPostTextMessageUnregisteredRecipient", id2, math.MaxInt64, content},
		{"testFailToPostTextMessageSenderEqualsRecipient", id1, id1, content},
		{"testFailToPostTextMessageHugeTextField", id2, id1, bigContent},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			resp, _ := requestPostMessage(c.senderid, c.recipientid, c.content, token2)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}
}

func testFailToPostImageMessageInvalidFields(t *testing.T) {
	dao.RefreshSchema(testDB)

	bigTextField := ""
	for i := 0; i < model.MAX_TEXT_FIELD_SIZE; i += 10 {
		bigTextField = bigTextField + "hugestring"
	}

	username1 := "user1"
	password1 := "pass1"

	username2 := "user2"
	password2 := "pass2"
	id1 := createUserSuccessfullyForTest(t, username1, password1)
	_ = createUserSuccessfullyForTest(t, username2, password2)
	id2, token2 := loginSuccessfullyForTest(t, username2, password2)
	content := model.ImageContent{"image", "url", 150, 120}
	bigUrl := model.ImageContent{"image", bigTextField, 150, 120}
	zeroHeight := model.ImageContent{"image", "url", 0, 120}
	zeroWidth := model.ImageContent{"image", "url", 150, 0}

	cases := []struct {
		name        string
		senderid    int64
		recipientid int64
		content     model.MessageContent
	}{
		{"testFailToPostImageMessageSenderZero", 0, id1, content},
		{"testFailToPostImageMessageRecipientZero", id2, 0, content},
		{"testFailToPostImageMessageUnregisteredSender", math.MaxInt64, id1, content},
		{"testFailToPostImageMessageUnregisteredRecipient", id2, math.MaxInt64, content},
		{"testFailToPostImageMessageSenderEqualsRecipient", id1, id1, content},
		{"testFailToPostImageMessageHugeUrlField", id2, id1, bigUrl},
		{"testFailToPostImageMessageZeroHeight", id2, id1, zeroHeight},
		{"testFailToPostImageMessageZeroWidth", id2, id1, zeroWidth},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			resp, _ := requestPostMessage(c.senderid, c.recipientid, c.content, token2)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}
}

func testFailToPostVideoMessageInvalidFields(t *testing.T) {
	dao.RefreshSchema(testDB)

	bigTextField := ""
	for i := 0; i < model.MAX_TEXT_FIELD_SIZE; i += 10 {
		bigTextField = bigTextField + "hugestring"
	}

	username1 := "user1"
	password1 := "pass1"

	username2 := "user2"
	password2 := "pass2"
	id1 := createUserSuccessfullyForTest(t, username1, password1)
	_ = createUserSuccessfullyForTest(t, username2, password2)
	id2, token2 := loginSuccessfullyForTest(t, username2, password2)
	content := model.VideoContent{"video", "youtube", "youtube"}
	bigUrl := model.VideoContent{"video", bigTextField, "youtube"}
	bigSource := model.VideoContent{"video", "youtube", bigTextField}
	invalid1 := model.VideoContent{"video", "youtube", "invalid"}
	invalid2 := model.VideoContent{"video", "invalid", "youtube"}

	cases := []struct {
		name        string
		senderid    int64
		recipientid int64
		content     model.MessageContent
	}{
		{"testFailToPostVideoMessageSenderZero", 0, id1, content},
		{"testFailToPostVideoMessageRecipientZero", id2, 0, content},
		{"testFailToPostVideoMessageUnregisteredSender", math.MaxInt64, id1, content},
		{"testFailToPostVideoMessageUnregisteredRecipient", id2, math.MaxInt64, content},
		{"testFailToPostVideoMessageSenderEqualsRecipient", id1, id1, content},
		{"testFailToPostVideoMessageHugeUrlField", id2, id1, bigUrl},
		{"testFailToPostVideoMessageHugeSourceField", id2, id1, bigSource},
		{"testFailToPostVideoMessageInvalidSourceField1", id2, id1, invalid1},
		{"testFailToPostVideoMessageInvalidSourceField2", id2, id1, invalid2},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			resp, _ := requestPostMessage(c.senderid, c.recipientid, c.content, token2)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}
}

func testFailToPostMessageInvalidBody(t *testing.T) {
	dao.RefreshSchema(testDB)

	requestReader := bytes.NewReader([]byte("invalid"))

	resp, _ := http.Post("http://localhost:8080/messages", "application/json", requestReader)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func testFailToPostMessageUnauthorized(t *testing.T) {
	dao.RefreshSchema(testDB)
	username1 := "user1"
	password1 := "pass1"

	username2 := "user2"
	password2 := "pass2"
	id1 := createUserSuccessfullyForTest(t, username1, password1)
	_ = createUserSuccessfullyForTest(t, username2, password2)
	id2, _ := loginSuccessfullyForTest(t, username2, password2)
	content := model.TextContent{"text", "text"}
	wrongToken, err := auth.GenerateToken()
	assert.Nil(t, err)
	req, _ := requestPostMessage(id2, id1, content, wrongToken.Uuid)
	assert.Equal(t, http.StatusUnauthorized, req.StatusCode)

}

func testFailToPostMessageWithClosedConnection(t *testing.T) {
	dao.RefreshSchema(testDB)

	token, _ := auth.GenerateToken()
	resp, err := requestPostMessage(int64(1), int64(2), model.TextContent{"text", "text"}, token.Uuid)
	if err == nil {
		defer resp.Body.Close()
	}
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func testFailToGetMessageInvalidQueryParams(t *testing.T) {
	dao.RefreshSchema(testDB)

	req, _ := http.NewRequest("GET", fmt.Sprintf("http://localhost:8080/messages?id=%v&start=%v&limit=%v", "invalid", 1, 1), nil)
	resp, _ := http.DefaultClient.Do(req)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	req, _ = http.NewRequest("GET", fmt.Sprintf("http://localhost:8080/messages?id=%v&start=%v&limit=%v", 1, "invalid", 1), nil)
	resp, _ = http.DefaultClient.Do(req)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	req, _ = http.NewRequest("GET", fmt.Sprintf("http://localhost:8080/messages?id=%v&start=%v&limit=%v", 1, 1, "invalid"), nil)
	resp, _ = http.DefaultClient.Do(req)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func testFailToGetMessageUnauthorized(t *testing.T) {
	dao.RefreshSchema(testDB)
	username1 := "user1"
	password1 := "pass1"

	username2 := "user2"
	password2 := "pass2"
	id1 := createUserSuccessfullyForTest(t, username1, password1)
	_ = createUserSuccessfullyForTest(t, username2, password2)
	id2, token := loginSuccessfullyForTest(t, username2, password2)
	content := model.TextContent{"text", "text"}
	messageID, _ := postMessageSuccessfullyForTest(t, id2, id1, content, token)
	wrongToken, err := auth.GenerateToken()
	assert.Nil(t, err)
	req, err := requestGetMessage(id1, messageID, 1, wrongToken.Uuid)
	assert.Equal(t, http.StatusUnauthorized, req.StatusCode)

}

func testFailToGetMessageInvalidFields(t *testing.T) {
	dao.RefreshSchema(testDB)

	username1 := "user1"
	password1 := "pass1"

	username2 := "user2"
	password2 := "pass2"
	id1 := createUserSuccessfullyForTest(t, username1, password1)
	_ = createUserSuccessfullyForTest(t, username2, password2)
	id2, token2 := loginSuccessfullyForTest(t, username2, password2)
	content := model.TextContent{"text", "text"}

	messageId, _ := postMessageSuccessfullyForTest(t, id2, id1, content, token2)

	cases := []struct {
		name  string
		id    int64
		start int64
		limit int64
	}{
		{"testFailToGetMessageRecipientZero", 0, messageId, 1},
		{"testFailToGetMessageStartZero", id2, 0, 1},
		{"testFailToGetMessageLimitZero", id2, messageId, 0},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			resp, _ := requestGetMessage(c.id, c.start, c.limit, token2)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}
}

func testFailToGetMessageWithClosedConnection(t *testing.T) {
	dao.RefreshSchema(testDB)

	token, _ := auth.GenerateToken()
	resp, err := requestGetMessage(int64(1), int64(1), 1, token.Uuid)
	if err == nil {
		defer resp.Body.Close()
	}
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func testGetMessages(t *testing.T) {
	dao.RefreshSchema(testDB)
	username1 := "user1"
	password1 := "pass1"

	username2 := "user2"
	password2 := "pass2"
	id1 := createUserSuccessfullyForTest(t, username1, password1)
	_ = createUserSuccessfullyForTest(t, username2, password2)
	_, token1 := loginSuccessfullyForTest(t, username1, password1)
	id2, token2 := loginSuccessfullyForTest(t, username2, password2)
	textContent := model.TextContent{"text", "text"}
	messageId1, stamp1 := postMessageSuccessfullyForTest(t, id2, id1, textContent, token2)

	imageContent := model.ImageContent{"image", "www.image.com", int64(150), int64(120)}
	messageId2, stamp2 := postMessageSuccessfullyForTest(t, id2, id1, imageContent, token2)

	videoContent := model.VideoContent{"video", "www.youtube.com", "youtube"}
	messageId3, stamp3 := postMessageSuccessfullyForTest(t, id2, id1, videoContent, token2)

	timestamp1, err := time.Parse("2006-01-02T15:04:05Z", stamp1)
	assert.Nil(t, err)

	textMessage := model.MessageResponse{messageId1, timestamp1, id2, id1, textContent}

	timestamp2, err := time.Parse("2006-01-02T15:04:05Z", stamp2)
	assert.Nil(t, err)

	imageMessage := model.MessageResponse{messageId2, timestamp2, id2, id1, imageContent}

	timestamp3, err := time.Parse("2006-01-02T15:04:05Z", stamp3)
	assert.Nil(t, err)

	videoMessage := model.MessageResponse{messageId3, timestamp3, id2, id1, videoContent}

	emptyMessages := []model.MessageResponse{}
	messages1 := []model.MessageResponse{textMessage, imageMessage, videoMessage}
	messages2 := []model.MessageResponse{imageMessage, videoMessage}
	messages3 := []model.MessageResponse{videoMessage}
	messages4 := []model.MessageResponse{textMessage, imageMessage}
	messages5 := []model.MessageResponse{imageMessage}
	messages6 := []model.MessageResponse{textMessage}

	cases := []struct {
		name             string
		start            int64
		limit            int64
		expectedMessages []model.MessageResponse
	}{
		{"testGetMessages0", messageId3 + 1, 1, emptyMessages},
		{"testGetMessages1", messageId1, 3, messages1},
		{"testGetMessages2", messageId2, 2, messages2},
		{"testGetMessages3", messageId3, 1, messages3},
		{"testGetMessages4", messageId1, 2, messages4},
		{"testGetMessages5", messageId2, 1, messages5},
		{"testGetMessages6", messageId1, 1, messages6},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			actualMessages := requestGetMessageSuccessfullyForTest(tt, id1, c.start, c.limit, token1)
			require.Equal(tt, len(c.expectedMessages), len(actualMessages))
			for i := range c.expectedMessages {
				expectedMessage := c.expectedMessages[i]
				actualMessage := actualMessages[i]
				assert.Equal(tt, expectedMessage.Id, actualMessage.Id)
				ts1 := expectedMessage.Timestamp.Format("2006-01-02T15:04:05Z")
				ts2 := actualMessage.Timestamp.Format("2006-01-02T15:04:05Z")
				assert.Equal(tt, ts1, ts2)
				assert.Equal(tt, expectedMessage.Sender, actualMessage.Sender)
				assert.Equal(tt, expectedMessage.Recipient, actualMessage.Recipient)
				actualContent, ok := actualMessage.Content.(map[string]interface{})
				assert.True(tt, ok)

				switch actualContent["type"] {
				case "text":
					expectedContent, ok := expectedMessage.Content.(model.TextContent)
					assert.True(tt, ok)
					assert.Equal(tt, expectedContent.Type, actualContent["type"])
					assert.Equal(tt, expectedContent.Text, actualContent["text"])
				case "image":
					expectedContent, ok := expectedMessage.Content.(model.ImageContent)
					assert.True(tt, ok)
					assert.Equal(tt, expectedContent.Url, actualContent["url"])
					height, ok := actualContent["height"].(float64)
					assert.True(tt, ok)
					assert.Equal(tt, expectedContent.Height, int64(height))
					width, ok := actualContent["width"].(float64)
					assert.True(tt, ok)
					assert.Equal(tt, expectedContent.Width, int64(width))
				case "video":
					expectedContent, ok := expectedMessage.Content.(model.VideoContent)
					assert.True(tt, ok)
					assert.Equal(tt, expectedContent.Url, actualContent["url"])
					assert.Equal(tt, expectedContent.Source, actualContent["source"])
				default:
					tt.Fail()
				}

			}
		})
	}
}

func createUserSuccessfullyForTest(t *testing.T, username string, password string) int64 {

	resp, err := requestCreateUser(username, password)
	require.Nil(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	require.Nil(t, err)

	response := model.UserResponseDTO{}
	json.Unmarshal(body, &response)
	return response.Id
}

func requestCreateUser(username string, password string) (*http.Response, error) {
	request := model.UserRequestDTO{
		username,
		password,
	}
	requestByte, _ := json.Marshal(request)
	requestReader := bytes.NewReader(requestByte)

	return http.Post("http://localhost:8080/users", "application/json", requestReader)
}

func loginSuccessfullyForTest(t *testing.T, username string, password string) (int64, string) {

	resp, err := requestLogin(username, password)
	require.Nil(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	require.Nil(t, err)

	response := model.LoginResponseDTO{}
	json.Unmarshal(body, &response)
	return response.Id, response.Token
}

func requestLogin(username string, password string) (*http.Response, error) {
	request := model.UserRequestDTO{
		username,
		password,
	}
	requestByte, _ := json.Marshal(request)
	requestReader := bytes.NewReader(requestByte)

	return http.Post("http://localhost:8080/login", "application/json", requestReader)
}

func postMessageSuccessfullyForTest(t *testing.T, sender int64, recipient int64, content model.MessageContent, token string) (int64, string) {

	resp, err := requestPostMessage(sender, recipient, content, token)
	require.Nil(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	require.Nil(t, err)

	response := model.PostMessageResponseDTO{}
	json.Unmarshal(body, &response)
	return response.Id, response.Timestamp
}

func requestPostMessage(sender int64, recipient int64, content model.MessageContent, token string) (*http.Response, error) {
	request := model.PostMessageRequestDTO{
		sender,
		recipient,
		content,
	}
	requestByte, _ := json.Marshal(request)
	requestReader := bytes.NewReader(requestByte)
	req, err := http.NewRequest("POST", "http://localhost:8080/messages", requestReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	return http.DefaultClient.Do(req)
}

func requestGetMessageSuccessfullyForTest(t *testing.T, recipient int64, start int64, limit int64, token string) []model.MessageResponse {

	resp, err := requestGetMessage(recipient, start, limit, token)
	require.Nil(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	require.Nil(t, err)

	response := model.GetMessageResponseDTO{}
	json.Unmarshal(body, &response)
	return response.Messages
}

func requestGetMessage(recipient int64, start int64, limit int64, token string) (*http.Response, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:8080/messages?id=%v&start=%v&limit=%v", recipient, start, limit), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	return http.DefaultClient.Do(req)
}
