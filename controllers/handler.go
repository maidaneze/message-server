package controllers

import "github.com/maidaneze/message-server/dao"

type Handler struct {
	Db dao.DB
}
