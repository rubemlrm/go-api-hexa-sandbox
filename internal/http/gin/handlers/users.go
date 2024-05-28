package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/rubemlrm/go-api-bootstrap/user"
	"log"
	"net/http"
)

type Controller struct{}

func (s *server) AddUser(c *gin.Context) {
	var uc *user.UserCreate
	if err := c.ShouldBindJSON(&uc); err != nil {
		s.Logger.Error("user", "creation", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := s.UserService.Create(uc)
	if err != nil {
		s.Logger.Error("user", "creation", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
	return
}

func (s *server) ListUsers(c *gin.Context) {
	res, err := s.UserService.All()
	if err != nil {
		log.Fatalf("error %s", err)
	}
	log.Printf("result %s", res)
}

func (s *server) GetUser(c *gin.Context, userId int) {
	res, err := s.UserService.Get(user.ID(int(userId)))
	if err != nil {
		s.Logger.Error("user", "get", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
	return
}
