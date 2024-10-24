// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"context"
)

type Querier interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	GetMessages(ctx context.Context, arg GetMessagesParams) ([]Message, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserById(ctx context.Context, id int64) (User, error)
	InsertMessage(ctx context.Context, arg InsertMessageParams) (Message, error)
}

var _ Querier = (*Queries)(nil)
