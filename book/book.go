package book

import (
	"github.com/fyk7/go-fiber-tutorial/database"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
)

type Book struct {
	gorm.Model
	Title  string `json:"title"`
	Author string `json:"author"`
	Rating int    `json:"rating"`
}

func GetBooks(c *fiber.Ctx) error {
	db := database.DBConn
	var books []Book
	db.Find(&books)
	return c.JSON(books)
}

func GetBook(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var book Book
	db.Find(&book, id)
	return c.JSON(book)
}

func NewBook(c *fiber.Ctx) error {
	db := database.DBConn

	book := new(Book)
	// book構造体のポインタにリクエストボディの内容をパースする
	if err := c.BodyParser(book); err != nil {
		return c.Status(503).SendString("Failed to create new book.")
	}

	db.Create(&book)
	return c.JSON(book)
}

func DeleteBook(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var book Book
	db.First(&book, id)
	if book.Title == "" {
		return c.Status(500).SendString("No book found with given ID!")
	}
	db.Delete(&book)
	return c.SendString("Book successfully deleted!")

}
