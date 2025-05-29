package ports

import (
	"context"
	"fmt"
	gin_handler "github.com/rubemlrm/go-api-bootstrap/internal/common/http/gin"
	"log/slog"
	"net/http"
	"time"

	"github.com/rubemlrm/go-api-bootstrap/internal/user/app"
	"github.com/rubemlrm/go-api-bootstrap/internal/user/app/query"
	"github.com/rubemlrm/go-api-bootstrap/internal/user/domain/user"

	"go.opentelemetry.io/otel"

	"github.com/gin-gonic/gin"
)

var tracer = otel.Tracer("gin-server")

type HTTPServer struct {
	app    app.UserModule
	Logger *slog.Logger
}

func NewHTTPServer(application app.UserModule, l *slog.Logger) ServerInterface {
	return HTTPServer{
		app:    application,
		Logger: l,
	}
}

func (s HTTPServer) AddUser(c *gin.Context) {
	req := c.MustGet("userCreate").(*user.UserCreate)
	_, span := tracer.Start(c.Request.Context(), "AddUser")
	defer span.End()

	reqID := c.GetString("requestID")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	ctx = context.WithValue(ctx, "requestID", reqID)
	defer cancel()

	id, err := s.app.Commands.CreateUser.Handle(c, req)
	if err != nil {
		s.Logger.Error("user", "creation", "error", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (s HTTPServer) ListUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	reqID := c.GetString("requestID")
	ctx = context.WithValue(ctx, "requestID", reqID)
	defer cancel()
	res, err := s.app.Queries.GetUsers.Handle(ctx, query.UserSearchFilters{})
	if err != nil {
		s.Logger.Error("user", "list", "error", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (s HTTPServer) GetUser(c *gin.Context, userID int) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	reqID := c.GetString("requestID")
	ctx = context.WithValue(ctx, "requestID", reqID)
	defer cancel()
	res, err := s.app.Queries.GetUser.Handle(ctx, query.UserSearch{ID: user.ID(userID)})
	if err != nil {
		s.Logger.Error("user", "get", "error", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if res == nil {
		s.Logger.Warn("user", "get", "not found user", fmt.Sprintf("%b", userID), nil)
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}

func RegisterHandlersWithOptionsAndValidations(router gin.IRouter, si ServerInterface, options GinServerOptions, logger *slog.Logger) {
	errorHandler := options.ErrorHandler
	if errorHandler == nil {
		errorHandler = func(c *gin.Context, err error, statusCode int) {
			c.JSON(statusCode, gin.H{"msg": err.Error()})
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandler:       errorHandler,
	}

	router.GET(options.BaseURL+"/users", wrapper.ListUsers)
	router.POST(options.BaseURL+"/users", gin_handler.ValidateRequestBody[*user.UserCreate](logger, "userCreate"), wrapper.AddUser)
	router.GET(options.BaseURL+"/users/:userId", wrapper.GetUser)
}
