package main

import (
	"fmt"
	"log"

	"github.com/fyk7/go-fiber-tutorial/book"
	"github.com/fyk7/go-fiber-tutorial/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/websocket/v2"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func getUser(c *fiber.Ctx) error {
	return c.JSON(&User{"fyk7", 26})
}

func helloWorld(c *fiber.Ctx) error {
	return c.SendString("hello world")
}

func getName(c *fiber.Ctx) error {
	msg := fmt.Sprintf("hello, %s", c.Params("name"))
	return c.SendString(msg)
}

func getJson(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "hello, fiber!",
	})
}

func setupRoutes(app *fiber.App) {
	app.Get("/api/v1/book", book.GetBooks)
	app.Get("/api/v1/book/:id", book.GetBook)
	app.Post("/api/v1/book", book.NewBook)
	app.Delete("/api/v1/book/:id", book.DeleteBook)
}

func initDatabase() {
	var err error
	database.DBConn, err = gorm.Open("sqlite3", "books.db")
	if err != nil {
		panic("Failed to connect to database")
	}
	fmt.Println("Database connection was established!")

	database.DBConn.AutoMigrate(&book.Book{})
	fmt.Println("Dababase Migrated!")
}

func main() {
	app := fiber.New()
	initDatabase()
	defer database.DBConn.Close()

	app.Use(logger.New())
	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	setupRoutes(app)

	app.Get("/", helloWorld)
	app.Get("/user", getUser)
	app.Get("/json", getJson)
	app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
		// c.Locals is added to the *websocket.Conn
		log.Println(c.Locals("allowed"))  // true
		log.Println(c.Params("id"))       // 123
		log.Println(c.Query("v"))         // 1.0
		log.Println(c.Cookies("session")) // ""

		// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
		var (
			mt  int
			msg []byte
			err error
		)
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", msg)

			if err = c.WriteMessage(mt, msg); err != nil {
				log.Println("write:", err)
				break
			}
		}
		// ?????????????????????????????????????????????websocket?????????????????????????????????
		// https://www.websocket.org/echo.html
		// ????????????websocket server???????????????
		// Access the websocket server: ws://localhost:3000/ws/123?v=1.0

	}))
	app.Get("/:name", getName)

	log.Fatal(app.Listen(":3000"))

}
