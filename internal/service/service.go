package service

import (
	"context"

	"github.com/himanshu-holmes/hms/internal/model"
)

type AuthService interface {
	Login(ctx context.Context, req model.LoginRequest)(*model.LoginResponse,error)
	CreateUser(ctx context.Context, req model.UserCreateRequest)(*model.User,error)
}

type PatientService interface {
	RegisterPatient(ctx context.Context, req model.)
}