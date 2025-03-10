package main

import (
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type User struct {
	ID     int    `json:"id"`
	Name   string `json:"name" validate:"required"`
	Pan    string `json:"pan" validate:"required,pan"`
	Mobile string `json:"mobile" validate:"required,mobile"`
	Email  string `json:"email" validate:"required,email"`
}

var validate *validator.Validate

var users []User
var currentID = 1

func panValidator(fl validator.FieldLevel) bool {
	panRegex := `^[A-Z]{5}[0-9]{4}[A-Z]$`
	match, _ := regexp.MatchString(panRegex, fl.Field().String())
	return match
}

func mobileValidator(fl validator.FieldLevel) bool {
	mobileRegex := `^[0-9]{10}$`
	match, _ := regexp.MatchString(mobileRegex, fl.Field().String())
	return match
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		log.Printf("Request processed in %s\n", duration)
	}
}

func CreateUsers(c *gin.Context) {
	var input []User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format", "details": err.Error()})
		return
	}

	var validUsers []User
	var validationErrors []map[string]interface{}

	for _, user := range input {
		if err := validate.Struct(user); err != nil {
			validationErrors = append(validationErrors, gin.H{
				"user": user,
				"reason": err.Error(),
			})
			continue
		}

		user.ID = currentID
		currentID++
		users = append(users, user)
		validUsers = append(validUsers, user)
	}

	response := gin.H{
		"success_count": len(validUsers),
		"failed_count": len(validationErrors),
		"users_created": validUsers,
	}

	if len(validationErrors) > 0 {
		response["validation_errors"] = validationErrors
	}

	c.JSON(http.StatusOK, response)
}

func main() {
	r := gin.Default()
	r.Use(Logger())

	validate = validator.New()
	validate.RegisterValidation("pan", panValidator)
	validate.RegisterValidation("mobile", mobileValidator)

	r.POST("/users", CreateUsers)
	r.Run(":8080")

}