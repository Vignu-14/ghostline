package routes

import (
	"anonymous-communication/backend/internal/handlers"

	fiberws "github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func registerWebSocketRoutes(app *fiber.App, websocketHandler *handlers.WebSocketHandler) {
	app.Use("/ws/chat", websocketHandler.Upgrade)
	app.Get("/ws/chat", fiberws.New(websocketHandler.HandleConnection))
}
