package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	svc := &Service{}

	engine := gin.New()

	engine.POST("/users", JSON(svc.CreateUser))
	engine.GET("/users", JSON(svc.ListUsers))
	engine.GET("/users/:id", JSON(svc.GetUser))

	if err := engine.Run(); err != nil {
		log.Fatal(err)
	}
}

type CreateUserRequest struct {
	User *User  `json:"user"`
	Opt  string `form:"opt"`
}

type CreateUserResponse struct {
	User *User  `json:"user"`
	Opt  string `json:"opt"`
}

type ListUsersRequest struct {
	Name string `form:"name"`
}

type ListUsersResponse struct {
	Users []*User `json:"users"`
}

type GetUserRequest struct {
	ID string `uri:"id"`
}

type GetUserResponse struct {
	User *User `json:"user"`
}

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Service struct{}

func (s *Service) CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
	return &CreateUserResponse{
		User: req.User,
		Opt:  req.Opt,
	}, nil
}

func (s *Service) ListUsers(ctx context.Context, req *ListUsersRequest) (*ListUsersResponse, error) {
	return &ListUsersResponse{
		Users: []*User{
			{
				ID:    "1",
				Name:  req.Name,
				Email: "min@mail.example.com",
			},
		},
	}, nil
}

func (s *Service) GetUser(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error) {
	if req.ID == "404" {
		return nil, Errorf(http.StatusNotFound, "user not found")
	}

	return &GetUserResponse{
		User: &User{
			ID:    req.ID,
			Name:  "min",
			Email: "min@mail.example.com",
		},
	}, nil
}

func JSON[Request, Response any](f func(context.Context, *Request) (*Response, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req Request

		if err := ctx.ShouldBindUri(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := ctx.ShouldBindQuery(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// If no request body will return io.EOF
		if err := ctx.ShouldBindJSON(&req); err != nil && err != io.EOF {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := f(ctx, &req)
		if err != nil {
			var apiErr APIError
			if errors.As(err, &apiErr) {
				ctx.JSON(apiErr.GetStatus(), gin.H{"error": apiErr.Error()})
				return
			}

			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, resp)
	}
}

func Errorf(status int, format string, a ...any) error {
	return &apiError{
		error:  fmt.Errorf(format, a...),
		status: status,
	}
}

type APIError interface {
	error

	GetStatus() int
}

type apiError struct {
	error

	status int
}

var _ APIError = &apiError{}

func (e *apiError) GetStatus() int {
	return e.status
}
