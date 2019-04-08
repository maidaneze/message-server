package model

type UserRequestDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserResponseDTO struct {
	Id int64 `json:"id"`
}

type LoginResponseDTO struct {
	Id    int64  `json:"id"`
	Token string `json:"token"`
}

type PostMessageRequestDTO struct {
	Sender    int64          `json:"sender"`
	Recipient int64          `json:"recipient"`
	Content   MessageContent `json:"content"`
}

type PostMessageResponseDTO struct {
	Id        int64  `json:"id"`
	Timestamp string `json:"timestamp"`
}

type GetMessageResponseDTO struct {
	Messages []MessageResponse `json:"messages"`
}
