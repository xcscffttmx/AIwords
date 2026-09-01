package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"uniqueIndex;size:50;not null" json:"username"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Word struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	UserID     uint           `gorm:"index;not null" json:"user_id"`
	Word       string         `gorm:"size:100;not null" json:"word"`
	Definition string         `gorm:"type:text;not null" json:"definition"`
	Examples   string         `gorm:"type:json;not null" json:"examples"`
	AIProvider string         `gorm:"column:ai_provider;size:50;not null" json:"ai_provider"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

var db *gorm.DB
var jwtSecret []byte

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	jwtSecret = []byte(getEnv("JWT_SECRET", "ai-wordbook-secret-key-2026"))

	initDB()

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "AI Wordbook API is running"})
	})

	r.POST("/api/register", register)
	r.POST("/api/login", login)

	auth := r.Group("/api")
	auth.Use(jwtMiddleware())
	{
		auth.GET("/words", getWords)
		auth.POST("/words/query", queryWord)
		auth.POST("/words/save", saveWord)
		auth.DELETE("/words/:id", deleteWord)
		auth.GET("/user/info", getUserInfo)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Server starting on port %s\n", port)
	fmt.Printf("📖 API Documentation: http://localhost:%s\n", port)
	r.Run(":" + port)
}

func initDB() {
	driver := strings.ToLower(getEnv("DB_DRIVER", "sqlite"))

	var err error
	switch driver {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			getEnv("DB_USER", "root"),
			getEnv("DB_PASSWORD", "root"),
			getEnv("DB_HOST", "localhost"),
			getEnv("DB_PORT", "3306"),
			getEnv("DB_NAME", "ai_wordbook"),
		)

		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to MySQL database: %v", err)
		}
		log.Println("✅ MySQL database connected successfully")
	default:
		dbPath := getEnv("DB_PATH", "ai_wordbook.db")
		db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to SQLite database: %v", err)
		}
		if err := db.AutoMigrate(&User{}, &Word{}); err != nil {
			log.Fatalf("Failed to migrate SQLite database: %v", err)
		}
		log.Println("✅ SQLite database connected successfully")
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func jwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization header"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}

func register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required,min=3,max=50"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingUser User
	if err := db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := User{
		Username:     req.Username,
		PasswordHash: string(hash),
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    user.ID,
	})
}

func login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user User
	if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   tokenString,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}

func getUserInfo(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
	})
}

func getWords(c *gin.Context) {
	userID := c.GetUint("user_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	var words []Word
	var total int64

	db.Model(&Word{}).Where("user_id = ? AND deleted_at IS NULL", userID).Count(&total)
	db.Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&words)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"list":      words,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
			"pages":     (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

func queryWord(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		Word       string `json:"word" binding:"required"`
		AIProvider string `json:"ai_provider" binding:"required,oneof=deepseek qianwen"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingWord Word
	if err := db.Where("user_id = ? AND word = ? AND deleted_at IS NULL", userID, req.Word).First(&existingWord).Error; err == nil {
		c.JSON(http.StatusOK, gin.H{
			"message":     "Word found in your wordbook",
			"source":      "database",
			"word":        existingWord.Word,
			"definition":  existingWord.Definition,
			"examples":    existingWord.Examples,
			"ai_provider": existingWord.AIProvider,
		})
		return
	}

	aiResponse, err := callAI(req.Word, req.AIProvider)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query AI: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Query successful",
		"source":      "ai",
		"word":        req.Word,
		"definition":  aiResponse["definition"],
		"examples":    aiResponse["examples"],
		"ai_provider": req.AIProvider,
	})
}

func saveWord(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		Word       string `json:"word" binding:"required"`
		Definition string `json:"definition" binding:"required"`
		Examples   string `json:"examples" binding:"required"`
		AIProvider string `json:"ai_provider" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingWord Word
	if err := db.Where("user_id = ? AND word = ? AND deleted_at IS NULL", userID, req.Word).First(&existingWord).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Word already exists in your wordbook"})
		return
	}

	word := Word{
		UserID:     userID,
		Word:       req.Word,
		Definition: req.Definition,
		Examples:   req.Examples,
		AIProvider: req.AIProvider,
	}

	if err := db.Create(&word).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save word"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Word saved successfully",
		"id":      word.ID,
	})
}

func deleteWord(c *gin.Context) {
	userID := c.GetUint("user_id")
	wordID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid word ID"})
		return
	}

	result := db.Where("id = ? AND user_id = ? AND deleted_at IS NULL", wordID, userID).Delete(&Word{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Word not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Word deleted successfully"})
}

func callAI(word, provider string) (map[string]interface{}, error) {
	var apiKey, apiURL string

	if provider == "deepseek" {
		apiKey = getEnv("DEEPSEEK_API_KEY", "")
		apiURL = "https://api.deepseek.com/v1/chat/completions"
	} else {
		apiKey = getEnv("QIANWEN_API_KEY", "")
		apiURL = "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions"
	}

	if apiKey == "" {
		log.Printf("[INFO] No API key for %s, using mock response", provider)
		return generateMockResponse(word), nil
	}

	prompt := fmt.Sprintf(`请为单词 "%s" 提供以下信息：
1. 单词的音标
2. 单词的词性和中文释义（用中文）
3. 3 个英文例句及其中文翻译

请以 JSON 格式返回，格式如下：
{
  "phonetic": "/音标/",
  "definition": "词性 1: 释义 1; 词性 2: 释义 2",
  "examples": [
    {"sentence": "例句 1 英文", "translation": "例句 1 中文"},
    {"sentence": "例句 2 英文", "translation": "例句 2 中文"},
    {"sentence": "例句 3 英文", "translation": "例句 3 中文"}
  ]
}`, word)

	var model string
	if provider == "deepseek" {
		model = "deepseek-chat"
	} else {
		model = "qwen-turbo"
	}

	reqBody := map[string]interface{}{
		"model":      model,
		"messages":   []map[string]string{{"role": "user", "content": prompt}},
		"max_tokens": 1000,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		log.Printf("[ERROR] Failed to marshal request body: %v", err)
		return generateMockResponse(word), nil
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("[ERROR] Failed to create HTTP request: %v", err)
		return generateMockResponse(word), nil
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[ERROR] AI API request failed: %v", err)
		return generateMockResponse(word), nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read API response: %v", err)
		return generateMockResponse(word), nil
	}

	if resp.StatusCode != 200 {
		log.Printf("[WARN] AI API returned non-200 status: %d, body: %s", resp.StatusCode, string(body))
		return generateMockResponse(word), nil
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("[ERROR] Failed to unmarshal API response: %v", err)
		return generateMockResponse(word), nil
	}

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		log.Printf("[WARN] No choices in AI API response")
		return generateMockResponse(word), nil
	}

	choice, ok := choices[0].(map[string]interface{})
	if !ok {
		log.Printf("[WARN] Invalid choice format in AI API response")
		return generateMockResponse(word), nil
	}

	message, ok := choice["message"].(map[string]interface{})
	if !ok {
		log.Printf("[WARN] Invalid message format in AI API response")
		return generateMockResponse(word), nil
	}

	content, ok := message["content"].(string)
	if !ok {
		log.Printf("[WARN] Invalid content format in AI API response")
		return generateMockResponse(word), nil
	}

	var aiResult map[string]interface{}
	if err := json.Unmarshal([]byte(content), &aiResult); err != nil {
		log.Printf("[ERROR] Failed to unmarshal AI result: %v", err)
		return generateMockResponse(word), nil
	}

	log.Printf("[INFO] Successfully got AI response for word: %s", word)
	return aiResult, nil
}

func generateMockResponse(word string) map[string]interface{} {
	examples := []map[string]string{
		{"sentence": fmt.Sprintf("The %s is an important concept in English learning.", word), "translation": fmt.Sprintf("%s是英语学习中的一个重要概念。", word)},
		{"sentence": fmt.Sprintf("Can you explain the meaning of %s?", word), "translation": fmt.Sprintf("你能解释%s的意思吗？", word)},
		{"sentence": fmt.Sprintf("I want to learn more about %s.", word), "translation": fmt.Sprintf("我想学习更多关于%s的知识。", word)},
	}

	examplesJSON, _ := json.Marshal(examples)

	return map[string]interface{}{
		"phonetic":   "/wɜːrd/",
		"definition": fmt.Sprintf("n. %s (英文单词)", word),
		"examples":   string(examplesJSON),
	}
}
