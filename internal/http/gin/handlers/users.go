package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rubemlrm/go-api-bootstrap/user"
	"net/http"
)

type Controller struct{}

func (s *server) AddUser(c *gin.Context) {
	var uc *user.UserCreate
	if err := c.ShouldBindJSON(&uc); err != nil {
		s.Logger.Error("user", "creation", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind"})
		return
	}

	id, err := s.UserService.Create(uc)
	if err != nil {
		s.Logger.Error("user", "creation", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
	return
}

func (s *server) ListUsers(c *gin.Context) {
	res, err := s.UserService.All()
	if err != nil {
		s.Logger.Error("user", "list", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (s *server) GetUser(c *gin.Context, userId int) {
	res, err := s.UserService.Get(user.ID(userId))
	if err != nil {
		s.Logger.Error("user", "get", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if res == nil {
		s.Logger.Warn("user", "get", "not found user", fmt.Sprintf("%b", userId), nil)
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
	return
}
