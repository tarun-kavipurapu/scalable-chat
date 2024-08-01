package internal

import (
	"log"
	"net/http"
	"tarun-kavipurapu/test-go-chat/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func SetupRouter(server *Server) *gin.Engine {
	r := server.router

	hub := NewHub()
	go hub.Run()
	userhandler := handlers.NewUserHandler(server.store)
	users := r.Group("/users")
	{

		users.POST("/login", userhandler.Login)
		users.POST("/signup", userhandler.Signup)

	}
	r.GET("/ws/:userID", func(c *gin.Context) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		}

		// Reading userID from request parameter
		userID := c.Param("userID")

		// Upgrading the HTTP connection to a WebSocket connection
		connection, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println(err)
			return
		}

		// Call the handler to create a new WebSocket user
		CreateNewSocketUser(hub, connection, userID)
	})

	return r
}
