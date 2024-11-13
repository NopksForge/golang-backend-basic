package user

import "github.com/google/uuid"

type SaveUserRequest struct {
	UserEmail string `json:"userEmail" validate:"required,email"`
	UserName  string `json:"userName" validate:"required"`
}

type SaveUserPayload struct {
	UserEmail string
	UserName  string
}

type SaveUserDto struct {
	UserId    uuid.UUID
	UserEmail string
	UserName  string
	Msg       string
}

type UpdateUserRequest struct {
	UserId    string `uri:"userId" validate:"required,uuid"`
	UserEmail string `json:"userEmail" validate:"required,email"`
	UserName  string `json:"userName" validate:"required"`
}

type UpdateUserPayload struct {
	UserId    uuid.UUID
	UserEmail string
	UserName  string
}

type DeleteUserRequest struct {
	UserId string `uri:"userId" validate:"required,uuid"`
}

type GetUserRequest struct {
	UserId string `uri:"userId" validate:"required,uuid"`
}

type GetUserDto struct {
	UserId    uuid.UUID
	UserEmail string
	UserName  string
}

type GetUserResponse struct {
	UserId    string `json:"userId"`
	UserEmail string `json:"userEmail"`
	UserName  string `json:"userName"`
}
