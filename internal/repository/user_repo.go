package repository

import (
	"context"

	"github.com/himanshu-holmes/hms/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type userRepo struct {
	queries *db.Queries
}

func NewUserRepo(queries *db.Queries) UserRepository {
	return &userRepo{queries: queries}
}

func (r *userRepo) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
	return r.queries.CreateUser(ctx, arg)
}

func (r *userRepo) GetUserByID(ctx context.Context, id pgtype.UUID) (db.User, error) {
	return r.queries.GetUserByID(ctx, id)
}

func (r *userRepo) GetUserByUsername(ctx context.Context, username string) (db.User, error) {
	return r.queries.GetUserByUsername(ctx, username)
}

func (r *userRepo) ListUsers(ctx context.Context, arg db.ListUsersParams) ([]db.User, error) {
	return r.queries.ListUsers(ctx, arg)
}

func (r *userRepo) UpdateUser(ctx context.Context, arg db.UpdateUserParams) (db.User, error) {
	return r.queries.UpdateUser(ctx, arg)
}

func (r *userRepo) SetUserActiveStatus(ctx context.Context, arg db.SetUserActiveStatusParams) (db.User, error) {
	return r.queries.SetUserActiveStatus(ctx, arg)
}

func (r *userRepo) DeleteUser(ctx context.Context, id pgtype.UUID) error {
	return r.queries.DeleteUser(ctx, id)
}
