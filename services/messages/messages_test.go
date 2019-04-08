package messages

import (
	"challenge/model"
	"testing"
	"time"

	"fmt"
	"net/http"
	"net/url"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshallMessageContent(t *testing.T) {
	zero := int64(0)
	one := int64(1)
	two := int64(2)
	cases := []struct {
		name                string
		request             model.PostMessageRequestDTO
		expectErrror        bool
		expectedRecipientId int64
		expectedSenderId    int64
		expectedType        string
		expectedText        string
		expectedUrl         string
		expectedHeight      int64
		expectedWidth       int64
		expectedSource      string
	}{
		{"testUnmarshallMessageContentWithEmptyRequest", model.PostMessageRequestDTO{}, true, zero, zero, "", "", "", zero, zero, ""},
		{"testUnmarshallMessageContentWithEmptyRequestContent", model.PostMessageRequestDTO{0, 0, struct{}{}}, true, zero, zero, "", "", "", zero, zero, ""},
		{"testUnmarshallMessageContentWithEmptyTextMessage", model.PostMessageRequestDTO{0, 0, model.TextContent{"", ""}}, true, zero, zero, "", "", "", zero, zero, ""},
		{"testUnmarshallMessageContentWithEmptyImageMessage", model.PostMessageRequestDTO{0, 0, model.ImageContent{"", "", 0, 0}}, true, zero, zero, "", "", "", zero, zero, ""},
		{"testUnmarshallMessageContentWithEmptyVideoMessage", model.PostMessageRequestDTO{0, 0, model.VideoContent{"", "", ""}}, true, zero, zero, "", "", "", zero, zero, ""},
		{"testUnmarshallMessageContentWithNonEmptyTextMessage", model.PostMessageRequestDTO{1, 2, map[string]interface{}{"type": "text", "text": "someText"}}, false, two, one, "text", "someText", "", zero, zero, ""},
		{"testUnmarshallMessageContentWithNonEmptyImageMessage", model.PostMessageRequestDTO{1, 2, map[string]interface{}{"type": "image", "url": "url", "height": 150, "width": 120}}, false, two, one, "image", "", "url", int64(150), int64(120), ""},
		{"testUnmarshallMessageContentWithNonEmptyVideoMessage", model.PostMessageRequestDTO{1, 2, map[string]interface{}{"type": "video", "url": "url", "source": "source"}}, false, two, one, "video", "", "url", zero, zero, "source"},
	}
	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			dto, err := UnmarshallMessageContent(c.request)
			if c.expectErrror {
				assert.NotNil(tt, err)
			} else {
				assert.Nil(tt, err)
				assert.Equal(tt, zero, dto.MessageId)
				assert.Equal(tt, c.expectedRecipientId, dto.RecipientId)
				assert.Equal(tt, c.expectedSenderId, dto.SenderId)
				assert.NotEqual(tt, time.Time{}, dto.Timestamp)
				assert.Equal(tt, c.expectedSenderId, dto.SenderId)
				assert.Equal(tt, c.expectedType, dto.Type)
				assert.Equal(tt, c.expectedText, dto.Text)
				assert.Equal(tt, c.expectedUrl, dto.Url)
				assert.Equal(tt, c.expectedHeight, dto.Height)
				assert.Equal(tt, c.expectedWidth, dto.Width)
				assert.Equal(tt, c.expectedSource, dto.Source)
			}
		})
	}
}

func TestValidMessageDto(t *testing.T) {
	emptyTextMessage := model.MessageDTO{}
	emptyTextMessage.RecipientId = 1
	emptyTextMessage.SenderId = 2
	emptyTextMessage.Type = "text"

	nonEmptyTextMessage := model.MessageDTO{}
	nonEmptyTextMessage.RecipientId = 1
	nonEmptyTextMessage.SenderId = 2
	nonEmptyTextMessage.Type = "text"
	nonEmptyTextMessage.Text = "text"

	emptyImageMessage := model.MessageDTO{}
	emptyImageMessage.Type = "image"

	validImageMessage := model.MessageDTO{}
	validImageMessage.RecipientId = 1
	validImageMessage.SenderId = 2
	validImageMessage.Type = "image"
	validImageMessage.Height = 156
	validImageMessage.Width = 156
	validImageMessage.Url = "http://www.google.com"

	invalidImageMessage1 := validImageMessage
	invalidImageMessage1.Height = 0
	invalidImageMessage2 := validImageMessage
	invalidImageMessage2.Width = 0

	emptyVideoMessage := model.MessageDTO{}
	emptyVideoMessage.Type = "video"

	validVideoMessage1 := model.MessageDTO{}
	validVideoMessage1.RecipientId = 1
	validVideoMessage1.SenderId = 2
	validVideoMessage1.Type = "video"
	validVideoMessage1.Url = "http://www.youtube.com"
	validVideoMessage1.Source = "youtube"

	validVideoMessage2 := validVideoMessage1
	validVideoMessage2.Url = "http://www.vimeo.com"
	validVideoMessage2.Source = "vimeo"

	invalidVideoMessage1 := validVideoMessage1
	invalidVideoMessage1.Url = "invalid"

	invalidVideoMessage2 := validVideoMessage2
	invalidVideoMessage2.Url = "invalid"

	invalidVideoMessage3 := validVideoMessage1
	invalidVideoMessage3.Source = "invalid"

	invalidVideoMessage4 := validVideoMessage2
	invalidVideoMessage4.Source = "invalid"

	invalidSenderId1 := nonEmptyTextMessage
	invalidSenderId1.SenderId = 0
	invalidSenderId2 := validImageMessage
	invalidSenderId2.SenderId = 0
	invalidSenderId3 := validVideoMessage1
	invalidSenderId3.SenderId = 0

	invalidRecipientId1 := nonEmptyTextMessage
	invalidRecipientId1.RecipientId = 0
	invalidRecipientId2 := validImageMessage
	invalidRecipientId2.RecipientId = 0
	invalidRecipientId3 := validVideoMessage1
	invalidRecipientId3.RecipientId = 0

	invalidSenderAndRecipientId1 := nonEmptyTextMessage
	invalidSenderAndRecipientId1.RecipientId = 2
	invalidSenderAndRecipientId2 := validImageMessage
	invalidSenderAndRecipientId2.RecipientId = 2
	invalidSenderAndRecipientId3 := validVideoMessage1
	invalidSenderAndRecipientId3.RecipientId = 2

	hugeString := string(make([]byte, model.MAX_TEXT_FIELD_SIZE*10))

	invalidMessageFieldLength1 := nonEmptyTextMessage
	invalidMessageFieldLength1.Url = hugeString
	invalidMessageFieldLength2 := nonEmptyTextMessage
	invalidMessageFieldLength2.Source = hugeString
	invalidMessageFieldLength3 := nonEmptyTextMessage
	invalidMessageFieldLength3.Text = hugeString

	invalidMessageFieldLength4 := validImageMessage
	invalidMessageFieldLength4.Url = hugeString
	invalidMessageFieldLength5 := validImageMessage
	invalidMessageFieldLength5.Source = hugeString
	invalidMessageFieldLength6 := validImageMessage
	invalidMessageFieldLength6.Text = hugeString

	invalidMessageFieldLength7 := validVideoMessage1
	invalidMessageFieldLength7.Url = hugeString
	invalidMessageFieldLength8 := validVideoMessage1
	invalidMessageFieldLength8.Source = hugeString
	invalidMessageFieldLength9 := validVideoMessage1
	invalidMessageFieldLength9.Text = hugeString

	cases := []struct {
		name     string
		message  model.MessageDTO
		expected bool
	}{
		{"testValidMessageDtoEmptyMessage", model.MessageDTO{}, false},
		{"testValidMessageDtoEmptyTextMessage", emptyTextMessage, true},
		{"testValidMessageDtoNonEmptyTextMessage", nonEmptyTextMessage, true},
		{"testValidMessageDtoEmptyImageMessage", emptyImageMessage, false},
		{"testValidMessageDtoValidImageMessage", validImageMessage, true},
		{"testValidMessageDtoInvalidImageMessage1", invalidImageMessage1, false},
		{"testValidMessageDtoInvalidImageMessage2", invalidImageMessage2, false},
		{"testValidMessageDtoEmptyVideoMessage", emptyVideoMessage, false},
		{"testValidMessageDtoValidVideoMessage1", validVideoMessage1, true},
		{"testValidMessageDtoValidVideoMessage2", validVideoMessage2, true},
		{"testValidMessageDtoInvalidVideoMessage1", invalidVideoMessage1, false},
		{"testValidMessageDtoInvalidVideoMessage2", invalidVideoMessage2, false},
		{"testValidMessageDtoInvalidVideoMessage3", invalidVideoMessage3, false},
		{"testValidMessageDtoInvalidVideoMessage4", invalidVideoMessage4, false},
		{"testValidMessageDtoInvalidRecipientId1", invalidRecipientId1, false},
		{"testValidMessageDtoInvalidRecipientId2", invalidRecipientId2, false},
		{"testValidMessageDtoInvalidRecipientId3", invalidRecipientId3, false},
		{"testValidMessageDtoInvalidSenderId1", invalidSenderId1, false},
		{"testValidMessageDtoInvalidSenderId2", invalidSenderId2, false},
		{"testValidMessageDtoInvalidSenderId3", invalidSenderId3, false},
		{"testValidMessageDtoInvalidSenderAndRecipientId1", invalidSenderAndRecipientId1, false},
		{"testValidMessageDtoInvalidSenderAndRecipientId2", invalidSenderAndRecipientId2, false},
		{"testValidMessageDtoInvalidSenderAndRecipientId3", invalidSenderAndRecipientId3, false},
		{"testValidMessageDtoInvalidMessageFieldLength1", invalidMessageFieldLength1, false},
		{"testValidMessageDtoInvalidMessageFieldLength2", invalidMessageFieldLength2, false},
		{"testValidMessageDtoInvalidMessageFieldLength3", invalidMessageFieldLength3, false},
		{"testValidMessageDtoInvalidMessageFieldLength4", invalidMessageFieldLength4, false},
		{"testValidMessageDtoInvalidMessageFieldLength5", invalidMessageFieldLength5, false},
		{"testValidMessageDtoInvalidMessageFieldLength6", invalidMessageFieldLength6, false},
		{"testValidMessageDtoInvalidMessageFieldLength7", invalidMessageFieldLength7, false},
		{"testValidMessageDtoInvalidMessageFieldLength8", invalidMessageFieldLength8, false},
		{"testValidMessageDtoInvalidMessageFieldLength9", invalidMessageFieldLength9, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			actual := ValidMessageDto(c.message)
			assert.Equal(tt, c.expected, actual)

		})
	}
}

func TestParseGetMessageQueryParams(t *testing.T) {
	cases := []struct {
		name          string
		id            string
		start         string
		limit         string
		success       bool
		expectedId    int64
		expectedStart int64
		expectedLimit int64
	}{
		{"testParseGetMessageQueryParamsEmptyParams1", "", "", "", false, 0, 0, 0},
		{"testParseGetMessageQueryParamsEmptyParams2", "1", "1", "", false, 0, 0, 0},
		{"testParseGetMessageQueryParamsEmptyParams3", "1", "", "1", false, 0, 0, 0},
		{"testParseGetMessageQueryParamsEmptyParams4", "", "1", "1", false, 0, 0, 0},
		{"testParseGetMessageQueryParamsInvalidParams1", "invalid", "1", "1", false, 0, 0, 0},
		{"testParseGetMessageQueryParamsInvalidParams1", "1", "invalid", "1", false, 0, 0, 0},
		{"testParseGetMessageQueryParamsInvalidParams1", "1", "1", "invalid", false, 0, 0, 0},
		{"testParseGetMessageQueryParamsvalidParams1", "1", "2", "3", true, 1, 2, 3},
		{"testParseGetMessageQueryParamsvalidParams2", "4", "5", "6", true, 4, 5, 6},
	}
	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			req := &http.Request{}
			var err error
			req.URL, err = url.Parse(fmt.Sprintf("www.example.com?id=%s&start=%s&limit=%s", c.id, c.start, c.limit))
			assert.Nil(tt, err)
			id, start, limit, err := ParseGetMessageQueryParams(req)
			assert.Equal(tt, c.expectedId, id)
			assert.Equal(tt, c.expectedStart, start)
			assert.Equal(tt, c.expectedLimit, limit)
			assert.True(tt, c.success == (err == nil))
		})
	}
}

func TestParseMessages(t *testing.T) {
	textMessage := model.MessageDTO{1, 2, 3, time.Now(), "text", "text", "", 0, 0, ""}
	imageMessage := model.MessageDTO{2, 2, 4, time.Now(), "image", "", "www.image.com", 150, 120, ""}
	videoMessage := model.MessageDTO{3, 2, 5, time.Now(), "video", "", "www.youtube.com", 0, 0, "youtube"}
	invalidMessage := model.MessageDTO{1, 2, 3, time.Now(), "invalid", "", "", 0, 0, ""}
	messages1 := []model.MessageDTO{textMessage, imageMessage, videoMessage, invalidMessage}
	responseMessages := ParseMessages(messages1)
	assert.Equal(t, 3, len(responseMessages))

	textMessageResponse := responseMessages[0]

	assert.Equal(t, textMessage.MessageId, textMessageResponse.Id)
	assert.Equal(t, textMessage.Timestamp, textMessageResponse.Timestamp)
	assert.Equal(t, textMessage.RecipientId, textMessageResponse.Recipient)
	assert.Equal(t, textMessage.SenderId, textMessageResponse.Sender)
	textContent, matches := textMessageResponse.Content.(model.TextContent)
	assert.True(t, matches)
	assert.Equal(t, textMessage.Text, textContent.Text)
	assert.Equal(t, textMessage.Type, textContent.Type)

	imageMessageResponse := responseMessages[1]

	assert.Equal(t, imageMessage.MessageId, imageMessageResponse.Id)
	assert.Equal(t, imageMessage.Timestamp, imageMessageResponse.Timestamp)
	assert.Equal(t, imageMessage.RecipientId, imageMessageResponse.Recipient)
	assert.Equal(t, imageMessage.SenderId, imageMessageResponse.Sender)
	imageContent, matches := imageMessageResponse.Content.(model.ImageContent)
	assert.True(t, matches)
	assert.Equal(t, imageMessage.Type, imageContent.Type)
	assert.Equal(t, imageMessage.Url, imageContent.Url)
	assert.Equal(t, imageMessage.Height, imageContent.Height)
	assert.Equal(t, imageMessage.Width, imageContent.Width)

	videoMessageResponse := responseMessages[2]

	assert.Equal(t, videoMessage.MessageId, videoMessageResponse.Id)
	assert.Equal(t, videoMessage.Timestamp, videoMessageResponse.Timestamp)
	assert.Equal(t, videoMessage.RecipientId, videoMessageResponse.Recipient)
	assert.Equal(t, videoMessage.SenderId, videoMessageResponse.Sender)
	videoContent, matches := videoMessageResponse.Content.(model.VideoContent)
	assert.True(t, matches)
	assert.Equal(t, videoMessage.Type, videoContent.Type)
	assert.Equal(t, videoMessage.Url, videoContent.Url)
	assert.Equal(t, videoMessage.Source, videoContent.Source)
}
