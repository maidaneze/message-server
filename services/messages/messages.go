package messages

import (
	"github.com/maidaneze/message-server/model"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var DEFAULT_LIMIT int64 = 100

//Unmarshals the PostMessageRequestDTO
//Returns error if the dto content isn't a valid TextMessage, ImageMessage or VideoMessage
//Returns the MessageDto and nil otherwise

func UnmarshallMessageContent(request model.PostMessageRequestDTO) (model.MessageDTO, error) {
	message := model.MessageDTO{}
	invalidBodyMessage := errors.New("Invalid body")
	content, ok := request.Content.(map[string]interface{})

	if !ok {
		return message, invalidBodyMessage
	}

	t, ok := content["type"].(string)

	if !ok {
		return message, invalidBodyMessage
	}

	switch t {
	case "text":
		text, ok := content["text"].(string)
		if !ok {
			return message, invalidBodyMessage
		}
		message.Text = text
		goto finish
	case "image":
		url, ok := content["url"].(string)
		if !ok {
			return message, invalidBodyMessage
		}
		message.Url = url

		switch v := content["height"].(type) {
		case float64:
			message.Height = int64(v)
		case float32:
			message.Height = int64(v)
		case int:
			message.Height = int64(v)
		default:
			return message, invalidBodyMessage
		}

		switch v := content["width"].(type) {
		case float64:
			message.Width = int64(v)
		case float32:
			message.Width = int64(v)
		case int:
			message.Width = int64(v)
		default:
			return message, invalidBodyMessage
		}
		goto finish
	case "video":
		url, ok := content["url"].(string)
		if !ok {
			return message, invalidBodyMessage
		}
		message.Url = url

		source, ok := content["source"].(string)
		if !ok {
			return message, invalidBodyMessage
		}
		message.Source = source
	default:
	}
finish:
	message.Type = t
	message.RecipientId = request.Recipient
	message.SenderId = request.Sender
	message.Timestamp = time.Now()
	return message, nil
}

//Validates if the messageDTO is a valid TextMessage, ImageMessage or VideoMessage
//Returns true if it is, false if it isn't
//The senderId and recipientID must be different and be greater than zero
//Maximum length for text fields is 1024 characters
//The type field must be either "text", "image" or "video"
//If the type field is "image" or "video", the url must be a valid uri
//If the type field is "image" the height and width must be greater than 0
//If the type field is "vidoe" the source must be youtube or vimeo

func ValidMessageDto(dto model.MessageDTO) bool {
	//Validate SenderId and RecipientID
	if dto.SenderId == dto.RecipientId || dto.SenderId <= 0 || dto.RecipientId <= 0 {
		return false
	}
	if len(dto.Url) > model.MAX_TEXT_FIELD_SIZE || len(dto.Source) > model.MAX_TEXT_FIELD_SIZE ||
		len(dto.Text) > model.MAX_TEXT_FIELD_SIZE {
		return false
	}

	switch dto.Type {
	//Validate TextMessage
	case "text":
		return true
	//Validate ImageMessage
	case "image":
		return dto.Height > 0 && dto.Width > 0
	//Validate VideoMessage
	case "video":
		return (dto.Source == "youtube" || dto.Source == "vimeo") && strings.Contains(dto.Url, dto.Source)
	default:
		return false
	}
	return true
}

//Parses the id, start and limit queryparams into int64
//Returns error if it cant

func ParseGetMessageQueryParams(r *http.Request) (int64, int64, int64, error) {
	//Get Query params
	recipient, foundRecipientId := r.URL.Query()["id"]
	message, foundMessageId := r.URL.Query()["start"]
	queyLimit, foundLimit := r.URL.Query()["limit"]

	if !(foundRecipientId && foundMessageId) {
		return 0, 0, 0, errors.New("Invalid Body")
	}

	recipientId, err := strconv.ParseInt(recipient[0], 10, 64)

	if err != nil {
		return 0, 0, 0, errors.New("Invalid Body")
	}

	messageId, err := strconv.ParseInt(message[0], 10, 64)

	if err != nil {
		return 0, 0, 0, errors.New("Invalid Body")
	}

	var limit int64 = DEFAULT_LIMIT

	if foundLimit {
		limit, err = strconv.ParseInt(queyLimit[0], 10, 64)

		if err != nil {
			return 0, 0, 0, errors.New("Invalid Body")
		}
	}
	return recipientId, messageId, limit, nil
}

//Parses the messages from the database into either textMessages, imageMessages or videoMessages
//Discards the invalid messages

func ParseMessages(m []model.MessageDTO) []model.MessageResponse {
	response := make([]model.MessageResponse, 0)
	for _, message := range m {
		messageResponse := model.MessageResponse{}

		messageResponse.Id = message.MessageId
		messageResponse.Timestamp = message.Timestamp
		messageResponse.Sender = message.SenderId
		messageResponse.Recipient = message.RecipientId
		switch message.Type {
		case "text":
			textContent := model.TextContent{message.Type, message.Text}
			messageResponse.Content = textContent
		case "image":
			imageContent := model.ImageContent{message.Type, message.Url, message.Height, message.Width}
			messageResponse.Content = imageContent
		case "video":
			videoContent := model.VideoContent{message.Type, message.Url, message.Source}
			messageResponse.Content = videoContent
		default:
			continue
		}
		response = append(response, messageResponse)
	}
	return response
}

func GetPostMessageResponseDTO(dto model.MessageDTO) model.PostMessageResponseDTO {
	return model.PostMessageResponseDTO{dto.MessageId, dto.Timestamp.Format("2006-01-02T15:04:05Z")}
}
