package internal

import (
	"context"
	"log"
	"net/http"
	"tarun-kavipurapu/test-go-chat/internal/handlers"
	"tarun-kavipurapu/test-go-chat/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func SetupRouter(server *Server) *gin.Engine {
	r := server.router

	userHandler := handlers.NewUserHandler(server.store)
	chatHandler := handlers.NewChatHandler(server.store)
	hub := NewHub(chatHandler, server.store)
	ctx := context.Background()
	go hub.Run(ctx)
	users := r.Group("/users")
	{

		users.POST("/login", userHandler.Login)
		users.POST("/signup", userHandler.Signup)

	}
	r.GET("/ws", func(c *gin.Context) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		}

		// Reading userID from request parameter
		token := c.Query("token")

		userId, err := utils.ExtractUserIdFromToken(token)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User Error Please Signout and login "})
			return
		}

		log.Println("userID,", userId)

		user, err := server.store.GetUserById(c, int64(userId))
		// log.Println(err.Error(), user)
		if err != nil {

			log.Println(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

			return
		}
		connection, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println(err)
			return
		}

		// Upgrading the HTTP connection to a WebSocket connection
		// Call the handler to create a new WebSocket user
		CreateNewSocketUser(hub, connection, user.ID)
	})

	return r
}
