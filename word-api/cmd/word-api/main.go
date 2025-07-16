package main

import (
	"fmt"
	"hash/fnv"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Word struct {
	gorm.Model
	Text   string `json:"text" gorm:"not null"`
	Color  string `json:"color" gorm:"not null"`
	Votes  int    `json:"votes" gorm:"default:0"`
	UserIP string `json:"user_ip" gorm:"not null"`
}

type Vote struct {
	gorm.Model
	WordID uint   `json:"word_id" gorm:"index"`
	UserIP string `json:"user_ip" gorm:"not null"`
}

var db *gorm.DB

func main() {
	// Инициализация базы данных
	initDB()

	r := gin.Default()

	// Middleware для логирования
	r.Use(loggingMiddleware())

	// API endpoints
	api := r.Group("/api")
	{
		api.POST("/words", addWord)
		api.GET("/words", getWords)
		api.POST("/words/:id/vote", addVote)
	}

	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}

func initDB() {
	var err error
	db, err = gorm.Open(sqlite.Open("words.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Автомиграция таблиц
	err = db.AutoMigrate(&Word{}, &Vote{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
}

func loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("%s %s", c.Request.Method, c.Request.URL)
		c.Next()
	}
}

func addWord(c *gin.Context) {
	var input struct {
		Text string `json:"text" binding:"required,max=100"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Генерация цвета на основе IP
	userIP := c.ClientIP()
	color := generateColor(userIP)

	word := Word{
		Text:   input.Text,
		Color:  color,
		UserIP: userIP,
		Votes:  0,
	}

	result := db.Create(&word)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save word"})
		return
	}

	c.JSON(http.StatusCreated, word)
}

func getWords(c *gin.Context) {
	var words []Word
	result := db.Order("created_at desc").Find(&words)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch words"})
		return
	}

	c.JSON(http.StatusOK, words)
}

func addVote(c *gin.Context) {
	wordID := c.Param("id")
	userIP := c.ClientIP()

	// Проверяем, голосовал ли уже пользователь
	var existingVote Vote
	result := db.Where("word_id = ? AND user_ip = ?", wordID, userIP).First(&existingVote)
	if result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You already voted for this word"})
		return
	}

	// Создаем новый голос
	vote := Vote{
		WordID: wordID,
		UserIP: userIP,
	}

	// Обновляем счетчик голосов
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&vote).Error; err != nil {
			return err
		}
		return tx.Model(&Word{}).Where("id = ?", wordID).Update("votes", gorm.Expr("votes + 1")).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process vote"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "vote added"})
}

func generateColor(ip string) string {
	hash := fnv.New32a()
	hash.Write([]byte(ip))
	hashValue := hash.Sum32()

	// Генерация более насыщенных цветов
	r := uint8((hashValue & 0xFF0000) >> 16)
	g := uint8((hashValue & 0x00FF00) >> 8)
	b := uint8(hashValue & 0x0000FF)

	// Делаем цвета более яркими
	r = r/2 + 128
	g = g/2 + 128
	b = b/2 + 128

	return fmt.Sprintf("#%02X%02X%02X", r, g, b)
}
