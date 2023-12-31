package token

import (
	"database/sql"
	"context"
	"net/http"
	"time"
	"strconv"

	"github.com/sotirismorf/microservice/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	queries *database.Queries
}

type credentials struct {
	Username    string `json:"username,omitempty" binding:"required,max=64"`
	Password    string `json:"password,omitempty" binding:"required,max=64"`
}

func NewServiceToken(queries *database.Queries) *Service {
	return &Service{queries: queries}
}

func (s *Service) RegisterHandlers(router *gin.Engine) {
	router.POST("/token", s.GetToken)
}

func GenerateToken(user_id int64) (string, error) {

	token_lifespan,err := strconv.Atoi("1")

	if err != nil {
		return "",err
	}

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = user_id
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(token_lifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte("this_is_a_secret"))

}

func (s *Service) GetToken(c *gin.Context) {
	var request credentials

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password),bcrypt.DefaultCost)
	//
	// if err != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	// 	return
	// }

	user, err := s.queries.GetUser(context.Background(), request.Username)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))

	if err != nil {
		if  err == bcrypt.ErrMismatchedHashAndPassword {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	token, err := GenerateToken(user.ID)

	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	// Build response
	// c.IndentedJSON(http.StatusCreated, response)
	c.JSON(http.StatusOK, gin.H{"token": token})
}
