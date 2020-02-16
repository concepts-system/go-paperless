package web

import (
	"github.com/labstack/echo/v4"

	"github.com/concepts-system/go-paperless/domain"
)

type userResponse struct {
	Username string `json:"username"`
	Surname  string `json:"surname"`
	Forename string `json:"forename"`
	IsAdmin  bool   `json:"isAdmin"`
	IsActive bool   `json:"isActive"`
}

type (
	userSerializer struct {
		C echo.Context
		*domain.User
	}

	userListSerializer struct {
		C     echo.Context
		Users []domain.User
	}
)

// Response returns the API response for a given user.
func (s userSerializer) Response() userResponse {
	return userResponse{
		Username: string(s.Username),
		Surname:  string(s.Surname),
		Forename: string(s.Forename),
		IsAdmin:  s.IsAdmin,
		IsActive: s.IsActive,
	}
}

// Response returns the API response for a list of users.
func (s userListSerializer) Response() []userResponse {
	response := make([]userResponse, len(s.Users))

	for i, user := range s.Users {
		serializer := userSerializer{s.C, &user}
		response[i] = serializer.Response()
	}

	return response
}
