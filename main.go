package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Book struct {
	ID     int    `json:"id"`
	Book   string `json:"book"`
	Author string `json:"author"`
}

type Author struct {
	Author string `json:"author"`
	Count  int    `json:"count"`
}

func main() {
	db, err := gorm.Open("mysql", "root:Flarecrazy@123@(127.0.0.1:3306)/onito_golang?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	r := gin.Default()

	r.GET("/get-books", func(c *gin.Context) {
		var books []Book
		if err := db.Find(&books).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, books)
	})

	r.POST("/add-book", func(c *gin.Context) {
		var lastBook Book
		db.Last(&lastBook)

		var book Book
		book.ID = lastBook.ID + 1
		if err := c.ShouldBindJSON(&book); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if err := db.Create(&book).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Book added successfully!"})
	})

	r.PUT("/update-book/:id", func(c *gin.Context) {
		var book Book
		if err := c.ShouldBindJSON(&book); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		id := c.Param("id")

		if err := db.Model(&book).Where("id = ?", id).Updates(book).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "Book updated successfully!"})
	})

	r.GET("/get-authors", func(c *gin.Context) {
		var authors []Author

		if err := db.Table("books").Select("author, count(*) as count").Group("author").Scan(&authors).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, authors)
	})

	r.Run("localhost:5000")

}
