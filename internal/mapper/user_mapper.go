package mapper

import (
	"github.com/google/uuid"
	"github.com/himanshu-holmes/hms/internal/db"
	"github.com/himanshu-holmes/hms/internal/model"
	
)

func ConvertDBUserToModel(user db.User) model.User {
	return model.User{
		ID: func() uuid.UUID {
			if user.ID.Valid {
				return user.ID.Bytes
			}
			return uuid.UUID{}
		}(),
		Username: user.Username,
		Role:     model.UserRole(user.Role),
		FirstName: func() *string {
			if user.FirstName.Valid {
				return &user.FirstName.String
			}
			return nil
		}(),
		LastName: func() *string {
			if user.LastName.Valid {
				return &user.LastName.String
			}
			return nil
		}(),
		Email: func() *string {
			if user.Email.Valid {
				return &user.Email.String
			}
			return nil
		}(),
		IsActive:     user.IsActive.Bool,
		CreatedAt:    user.CreatedAt.Time,
		UpdatedAt:    user.UpdatedAt.Time,
		PasswordHash: "",
	}
}
