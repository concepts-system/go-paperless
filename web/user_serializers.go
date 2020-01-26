package web

import (
	"github.com/labstack/echo/v4"

	"github.com/concepts-system/go-paperless/domain"
)

// userResponse defines the user model projection returned by API methods.
type userResponse struct {
	Username string `json:"username"`
	Surname  string `json:"surname"`
	Forename string `json:"forename"`
	IsAdmin  bool   `json:"isAdmin"`
	IsActive bool   `json:"isActive"`
}

// userSerializer defines functionality for serializing user models into user responses.
type (
	userSerializer struct {
		C echo.Context
		*domain.User
	}

	// userListSerializer defines functionality for serializing users models into users responses.
	userListSerializer struct {
		C     echo.Context
		Users []domain.User
	}
)

// Response returns the API response for a given user model.
func (s userSerializer) Response() userResponse {
	return userResponse{
		Username: string(s.Username),
		Surname:  string(s.Surname),
		Forename: string(s.Forename),
		IsAdmin:  s.IsAdmin,
		IsActive: s.IsActive,
	}
}

// Response returns the API response for a given users model.
func (s userListSerializer) Response() []userResponse {
	response := make([]userResponse, len(s.Users))

	for idx, user := range s.Users {
		serializer := userSerializer{s.C, &user}
		response[idx] = serializer.Response()
	}

	return response
}
