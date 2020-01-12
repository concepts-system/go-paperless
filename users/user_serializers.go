package users

import (
	"time"

	"github.com/labstack/echo/v4"
)

// UserResponse defines the user model projection returned by API methods.
type UserResponse struct {
	ID        uint   `json:"id,omitempty"`
	Username  string `json:"username"`
	Surname   string `json:"surname"`
	Forename  string `json:"forename"`
	IsAdmin   bool   `json:"isAdmin"`
	IsActive  bool   `json:"isActive"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// UserSerializer defines functionality for serializing user models into user responses.
type (
	UserSerializer struct {
		C echo.Context
		*UserModel
	}

	// UserListSerializer defines functionality for serializing users models into users responses.
	UserListSerializer struct {
		C     echo.Context
		Users []UserModel
	}
)

// Response returns the API response for a given user model.
func (s *UserSerializer) Response() UserResponse {
	return UserResponse{
		ID:        s.ID,
		Username:  s.Username,
		Surname:   s.Surname,
		Forename:  s.Forename,
		IsAdmin:   s.IsAdmin,
		IsActive:  s.IsActive,
		CreatedAt: s.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: s.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

// Response returns the API response for a given users model.
func (s *UserListSerializer) Response() []UserResponse {
	response := make([]UserResponse, len(s.Users))
	for idx, user := range s.Users {
		serializer := UserSerializer{s.C, &user}
		response[idx] = serializer.Response()
	}

	return response
}
