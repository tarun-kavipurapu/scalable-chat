package handlers

import (
	"context"
	db "tarun-kavipurapu/test-go-chat/db/sqlc"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	store db.Store
}

func NewChatHandler(store db.Store) *ChatHandler {
	return &ChatHandler{store: store}
}

func (c *ChatHandler) InsertMessage(ctx context.Context, from int64, to int64, content string) error {

	_, err := c.store.InsertMessage(ctx, db.InsertMessageParams{
		FromUserID: from,
		ToUserID:   to,
		IsSent:     true, //change this IsSent accordingly and make necessary changes in future for adding any new features
		Content:    content,
	})

	if err != nil {
		return err
	}

	return nil
}

func (c *ChatHandler) GetMessage(ctx *gin.Context) {
	//get Messages By either sending as a Requst Object or a  Query Parameters
	// c.store.GetMessages(ctx, db.GetMessagesParams{
	// 	FromUserID: ,
	// })
}
